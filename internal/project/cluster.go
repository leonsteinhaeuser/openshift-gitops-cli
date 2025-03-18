package project

import (
	"fmt"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
)

var (
	_ AddonHandler = &Cluster{}
)

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
}

// DisableAddon disables the addon for the cluster by setting the enabled flag to false
func (c *Cluster) DisableAddon(addon string) {
	if _, ok := c.Addons[addon]; !ok {
		// already disabled or not found
		return
	}
	c.Addons[addon].Enabled = false
}

// GetAddons returns the cluster addons
func (c *Cluster) GetAddons() ClusterAddons {
	return c.Addons
}

// GetAddon returns the addon by name
func (c *Cluster) GetAddon(name string) *ClusterAddon {
	return c.Addons[name]
}

// Render renders the cluster configuration using the given project templates
func (c *Cluster) Render(config *ProjectConfig, env, stage string) error {
	properties := utils.MergeMaps(config.EnvStageProperty(env, stage), c.Properties)

	templates, err := template.LoadTemplateManifest(config.TemplateBasePath)
	if err != nil {
		return fmt.Errorf("failed to load base templates: %w", err)
	}

	addonProperties := c.AddonProperties(config, env, stage)
	addons := map[string]template.AddonData{}
	for k, v := range addonProperties {
		addons[k] = template.AddonData{
			Enabled:     v.Enabled,
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
			Environment:       env,
			Stage:             stage,
			Cluster:           c.Name,
			ClusterProperties: properties,
			Properties:        addonValue.Properties,
		})
		if err != nil {
			return fmt.Errorf("failed to render addon: %s, Error: %w", addonName, err)
		}
	}
	return nil
}

// SetDefaultAddons sets the default addons for the cluster
func (c *Cluster) SetDefaultAddons(config *ProjectConfig) {
	for addonName, addon := range config.Addons {
		if !addon.DefaultEnabled {
			// skip disabled addons
			continue
		}

		_, ok := c.Addons[addonName]
		if ok {
			// was found, we respect the user setting
			continue
		}

		cAddon := &ClusterAddon{
			Enabled:    true,
			Properties: map[string]any{},
		}

		for key, property := range config.ParsedAddons[addonName].Properties {
			cAddon.Properties[key] = property.Default
		}
		c.Addons[addonName] = cAddon
	}
}

// AddonProperties returns the addon properties for the cluster merged with the environment and stage properties
func (c *Cluster) AddonProperties(config *ProjectConfig, env, stg string) map[string]*ClusterAddon {
	properties := c.Addons
	for addonName, addon := range c.Addons {
		if !addon.Enabled {
			// addon was disabled on the cluster level, we skip it
			continue
		}
		envAddonProps := map[string]any{}
		if env := config.GetEnvironment(env).GetAddon(addonName); env != nil {
			envAddonProps = env.Properties
		}
		stageAddonProps := map[string]any{}
		if stg := config.GetStage(env, stg).GetAddon(addonName); stg != nil {
			stageAddonProps = stg.Properties
		}
		addonProps := map[string]any{}
		for key, property := range config.ParsedAddons[addonName].Properties {
			addonProps[key] = property.Default
		}
		properties[addonName].Properties = utils.MergeMaps(addonProps, envAddonProps, stageAddonProps, c.GetAddon(addonName).Properties)
	}
	return properties
}
