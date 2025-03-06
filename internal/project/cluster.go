package project

import (
	"fmt"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
)

type ClusterAddon struct {
	Enabled    bool           `json:"enabled"`
	Properties map[string]any `json:"properties"`
}

// AllRequiredPropertiesSet checks if all required properties are set for the addon
func (ca *ClusterAddon) AllRequiredPropertiesSet(config *ProjectConfig, addonName string) error {
	for key, property := range config.ParsedAddons[addonName].Properties {
		if property.Required && ca.Properties[key] == nil {
			return fmt.Errorf("property for key %s is required", key)
		}
		_, err := config.ParsedAddons[addonName].Properties[key].ParseValue(ca.Properties[key])
		if err != nil {
			return fmt.Errorf("property for key %s is invalid: %w", key, err)
		}
	}
	return nil
}

type Cluster struct {
	Name       string                   `json:"-"`
	Addons     map[string]*ClusterAddon `json:"addons"`
	Properties map[string]string        `json:"properties"`
}

// IsAddonEnabled checks if the addon is enabled for the cluster
func (c Cluster) IsAddonEnabled(addon string) bool {
	_, ok := c.Addons[addon]
	if !ok {
		return false
	}
	return c.Addons[addon].Enabled
}

// EnableAddon enables the addon for the cluster by setting the enabled flag to true
func (c *Cluster) EnableAddon(addon string) {
	if c.Addons[addon] == nil {
		c.Addons[addon] = &ClusterAddon{
			Enabled: true,
		}
		return
	}
	c.Addons[addon].Enabled = true
	fmt.Println("Enabled addon option:", c.Addons[addon].Enabled)
}

// DisableAddon disables the addon for the cluster by setting the enabled flag to false and removing all properties
func (c *Cluster) DisableAddon(addon string) {
	if _, ok := c.Addons[addon]; !ok {
		// already disabled or not found
		return
	}
	c.Addons[addon].Enabled = false
	c.Addons[addon].Properties = map[string]any{}
}

// Render renders the cluster configuration using the given project templates
func (c *Cluster) Render(config *ProjectConfig, env, stage string) error {
	properties := utils.MergeMaps(config.EnvStageProperty(env, stage), c.Properties)

	templates, err := template.LoadTemplateManifest(config.TemplateBasePath)
	if err != nil {
		return fmt.Errorf("failed to load base templates: %w", err)
	}

	addons := map[string]template.AddonData{}
	for k, v := range c.Addons {
		if !v.Enabled {
			continue
		}
		addons[k] = template.AddonData{
			Annotations: config.ParsedAddons[k].Annotations,
			Properties:  v.Properties,
		}
	}

	// render templates
	for _, t := range templates {
		err = t.Render(config.BasePath, template.TemplateData{
			Environment: env,
			Stage:       stage,
			ClusterName: c.Name,
			Properties:  properties,
			Addons:      addons,
		})
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
	}

	// render addons
	for addonName, addonValue := range addons {
		atc, err := template.LoadTemplatesFromAddonManifest(config.ParsedAddons[addonName])
		if err != nil {
			return fmt.Errorf("failed to load addon %s templates: %w, value: %+v", addonName, err, config.ParsedAddons[addonName])
		}
		err = atc.Render(config.BasePath, template.AddonTemplateData{
			Environment: env,
			Stage:       stage,
			Cluster:     c.Name,
			Properties:  addonValue.Properties,
		})
		if err != nil {
			return fmt.Errorf("failed to render addon: %s, Error: %w", addonName, err)
		}
	}
	return nil
}

func (c *Cluster) AllRequiredPropertiesSet(config *ProjectConfig) error {
	for addonName, addon := range c.Addons {
		if !addon.Enabled {
			// skip disabled addons
			continue
		}
		err := addon.AllRequiredPropertiesSet(config, addonName)
		if err != nil {
			return fmt.Errorf("failed to validate properties for addon %s: %w", addonName, err)
		}
	}
	return nil
}
