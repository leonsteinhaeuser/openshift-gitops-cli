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

// templateManifestPropertiesMenu handles the selection of properties for a single manifest.yaml file
func templateManifestPropertiesMenu(writer io.Writer, reader *bufio.Reader, c project.Cluster, tmpl template.TemplateManifest) (map[string]any, error) {
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

	for {
		addonProps := utils.MapKeysToList(tmpl.Properties)
		if len(addonProps) == 0 {
			// no properties to set
			break
		}

		prompt := promptui.Select{
			Label:     "Property",
			Items:     addonProps,
			Templates: helperSelectTemplateAddonProperties(c, tmpl),
		}
		_, propertyKey, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		if propertyKey == "" {
			return nil, fmt.Errorf("property key cannot be empty")
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
	}
	return c.Addons[tmpl.Name], nil
}
