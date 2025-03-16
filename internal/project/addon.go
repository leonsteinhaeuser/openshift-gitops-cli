package project

import (
	"fmt"
)

type AddonHandler interface {
	// IsAddonEnabled checks if the addon is enabled
	IsAddonEnabled(name string) bool
	// EnableAddon enables the addon
	EnableAddon(name string)
	// DisableAddon disables the addon
	DisableAddon(name string)
	// GetAddons returns the addons
	GetAddons() ClusterAddons
	// GetAddon returns the addon by name
	GetAddon(name string) *ClusterAddon
}

type ClusterAddons map[string]*ClusterAddon

func (ca ClusterAddons) AllRequiredPropertiesSet(config *ProjectConfig) error {
	for addonName, addon := range ca {
		if !addon.Enabled {
			fmt.Printf("addon %s is disabled\n", addonName)
			continue
		}
		err := addon.AllRequiredPropertiesSet(config, addonName)
		if err != nil {
			return fmt.Errorf("failed to validate addon %s: %w", addonName, err)
		}
	}
	return nil
}

func (ca ClusterAddons) IsEnabled(addon string) bool {
	if ca[addon] == nil {
		return false
	}
	return ca[addon].IsEnabled()
}

type ClusterAddon struct {
	Enabled    bool           `json:"enabled"`
	Properties map[string]any `json:"properties"`
}

// AllRequiredPropertiesSet checks if all required properties are set for the addon
func (ca *ClusterAddon) AllRequiredPropertiesSet(config *ProjectConfig, addonName string) error {
	for key, property := range config.ParsedAddons[addonName].Properties {
		if property.Required && ca.Properties[key] == nil {
			return fmt.Errorf("[%s] property for key %s is required", addonName, key)
		}
		_, err := property.ParseValue(ca.Properties[key])
		if err != nil {
			return fmt.Errorf("[%s] property for key %s is invalid: %w", addonName, key, err)
		}
	}
	return nil
}

// IsEnabled checks if the addon is enabled
func (ca ClusterAddon) IsEnabled() bool {
	return ca.Enabled
}

func (ca *ClusterAddon) SetProperty(key string, value any) {
	if ca.Properties == nil {
		ca.Properties = map[string]any{}
	}
	ca.Properties[key] = value
}
