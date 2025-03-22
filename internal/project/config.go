package project

import (
	"errors"
	"fmt"
	"os"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"sigs.k8s.io/yaml"
)

// ParseConfig reads a yaml file from the given path and unmarshals it into a ProjectConfig struct
// If the config file does not exist, it will be created with default values
func ParseConfig(path string) (*ProjectConfig, error) {
	_, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to check for the PROJECT.yaml file: %w", err)
	}
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		_, err = f.WriteString("basePath: overlays/\ntemplateBasePath: templates/\n")
		if err != nil {
			return nil, err
		}
	}

	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pc := &ProjectConfig{}
	err = yaml.Unmarshal(bts, pc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config to ProjectConfig: %w", err)
	}

	if pc.Environments == nil {
		pc.Environments = map[string]*Environment{}
	}

	if pc.Addons == nil {
		pc.Addons = map[string]Addon{}
	}

	if pc.ParsedAddons == nil {
		pc.ParsedAddons = map[string]template.TemplateManifest{}
	}

	// load all addons, so we can use them later
	for k, v := range pc.Addons {
		tm, err := template.LoadManifest(v.Path)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while loading the addon [%s] manifest file: %s, %v", k, v.Path, err)
		}
		tm.Name = k
		tm.BasePath = v.Path
		tm.Group = v.Group
		pc.ParsedAddons[k] = *tm
	}
	return pc, nil
}

// UpdateOrCreateConfig writes a ProjectConfig struct to a yaml file at the given path
func UpdateOrCreateConfig(path string, config *ProjectConfig) error {
	bts, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal ProjectConfig to yaml: %w", err)
	}

	err = os.WriteFile(path, bts, 0644)
	if err != nil {
		return fmt.Errorf("failed to write ProjectConfig to file: %w", err)
	}
	return nil
}
