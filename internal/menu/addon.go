package menu

import (
	"bufio"
	"fmt"
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

func AddAddon(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) error {
	// ask for the addon name
	addonName, err := cli.StringQuestion(writer, reader, "Addon Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("addon name cannot be empty")
		}

		if _, ok := config.Addons[s]; ok {
			return fmt.Errorf("addon already exists")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// let's ask if the user want to add additional properties
	isEnabledByDefault, err := cli.BooleanQuestion(writer, reader, "Should this addon be enabled by default?", false)
	if err != nil {
		return err
	}

	sourcePath, err := cli.StringQuestion(writer, reader, "Please provide the path to the location of the addon (the directory must contain a manifest.yaml file)", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("addon source path cannot be empty")
		}
		return nil
	})
	if err != nil {
		return err
	}

	prompt := promptui.SelectWithAdd{
		Label:    "Select Group",
		Items:    config.AddonGroups(),
		AddLabel: "Create new group",
	}
	_, groupName, err := prompt.Run()
	if err != nil {
		return err
	}

	// let's ask if the user want to add additional properties
	confirmation, err := cli.BooleanQuestion(writer, reader, "Are you sure you want to create the addon?", false)
	if err != nil {
		return err
	}
	if !confirmation {
		return fmt.Errorf("canceled addon creation")
	}

	if config.Addons == nil {
		config.Addons = map[string]project.Addon{}
	}

	config.Addons[addonName] = project.Addon{
		DefaultEnabled: isEnabledByDefault,
		Path:           sourcePath,
		Group:          groupName,
	}
	return nil
}

const (
	cancelManageAddons = "Done"
	enableDisableAddon = "Enable / Disable"
)

func allRequiredPropertiesSet(addonConfig template.TemplateManifest, properties map[string]any) bool {
	for propKe, propVal := range addonConfig.Properties {
		if propVal.Required && properties[propKe] == nil {
			return false
		}
	}
	return true
}

func manageAddons(writer io.Writer, reader *bufio.Reader, config *project.ProjectConfig, env, stage, cluster string) (map[string]map[string]any, error) {
	addonConfig := map[string]map[string]any{}
OUTER:
	for {
		// provide a select menu to manage addons
		prompt := promptui.Select{
			Label:     "Enable / Disable Addons",
			Items:     append(utils.MapKeysToList(config.ParsedAddons), cancelManageAddons),
			Templates: helperManageAddons(config, addonConfig),
		}
		_, selectValue, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if selectValue == cancelManageAddons {
			// check if all required properties are set
			for addonName, addonSettings := range config.ParsedAddons {
				if !allRequiredPropertiesSet(addonSettings, addonConfig[addonName]) {
					// not all required properties are set
					fmt.Fprintf(writer, "You must set all required properties for addon %s to proceed!\n", addonName)
					continue OUTER
				}
			}
			// user is done with managing addons
			break
		}

		// ask for addon properties
		addonProperties, err := addonPropertiesMenu(
			writer,
			reader,
			config.Environments[env].Stages[stage].Clusters[cluster],
			config.ParsedAddons[selectValue])
		if err != nil {
			return nil, err
		}
		addonConfig[selectValue] = addonProperties
	}
	return addonConfig, nil
}

// addonPropertiesMenu handles the selection of properties for a single addon
func addonPropertiesMenu(writer io.Writer, reader *bufio.Reader, c project.Cluster, tmpl template.TemplateManifest) (map[string]any, error) {
	if c.Addons == nil {
		c.Addons = map[string]map[string]any{}
	}
	if c.Addons[tmpl.Name] == nil {
		c.Addons[tmpl.Name] = map[string]any{}
	}
	// add default values and required properties to the cluster
	for k, v := range tmpl.Properties {
		if v.Default != nil || v.Required {
			c.Addons[tmpl.Name][k] = v.Default
		}
	}

OUTER:
	for {
		addonProps := utils.MapKeysToList(tmpl.Properties)
		if len(addonProps) == 0 {
			// no properties to set
			break
		}

		prompt := promptui.Select{
			Label:     "Property",
			Items:     append([]string{enableDisableAddon}, append(addonProps, cancelManageAddons)...),
			Templates: helperSelectTemplateAddonProperties(c, tmpl),
		}
		_, propertyKey, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		switch propertyKey {
		case cancelManageAddons:
			if !allRequiredPropertiesSet(tmpl, c.Addons[tmpl.Name]) {
				// not all required properties are set
				fmt.Fprintf(writer, "You must set all required properties to proceed!\n")
				continue
			}

			// user is done with managing properties
			break OUTER
		case enableDisableAddon:
			// user wants to enable or disable the addon
			shouldAddonBeEnabled, err := cli.BooleanQuestion(writer, reader, fmt.Sprintf("Do you really want to enable %s?", tmpl.Name), c.Addons[tmpl.Name] != nil)
			if err != nil {
				return nil, err
			}
			if shouldAddonBeEnabled {
				if _, ok := c.Addons[tmpl.Name]; !ok {
					// addon is not enabled yet
					c.Addons[tmpl.Name] = map[string]any{}
				}
				// already enabled
			} else {
				delete(c.Addons, tmpl.Name)
			}
			continue
		}

		value, err := cli.UntypedQuestion(writer, reader, "Value", "", func(s any) error {
			if s == nil {
				return fmt.Errorf("value cannot be empty")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		// check if the value is of the correct type
		actualValue, err := tmpl.Properties[propertyKey].ParseValue(value)
		if err != nil {
			fmt.Fprintf(writer, "required property %s is not valid. %s\n", propertyKey, err.Error())
			continue
		}
		c.Addons[tmpl.Name][propertyKey] = actualValue

		/*
			confirmation, err := cli.BooleanQuestion(writer, reader, "Do you want to add or update another property?", false)
			if err != nil {
				return nil, err
			}
			if confirmation {
				// user wants to add or update another property
				continue
			}

			// check if all required properties are set
			allRequiredSet := true
			for k, v := range tmpl.Properties {
				if v.Required && c.Addons[tmpl.Name][k] == nil {
					allRequiredSet = false
					fmt.Fprintf(writer, "required property %s is missing\nYou must set all properties to proceed!\n", k)
				}
			}

			if allRequiredSet {
				// user has set all required properties and confirmed to be done
				break
			}
		*/
	}
	return c.Addons[tmpl.Name], nil
}
