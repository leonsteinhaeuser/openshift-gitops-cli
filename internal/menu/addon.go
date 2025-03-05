package menu

import (
	"bufio"
	"fmt"
	"io"

	"github.com/k0kubun/pp/v3"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

type addonClusterMenu struct {
	writer io.Writer
	reader *bufio.Reader
	config *project.ProjectConfig
}

func (a *addonClusterMenu) menuManageAddons(cluster *project.Cluster) error {
	for {
		prompt := promptui.Select{
			Label: "Manage Addons",
			Items: append(utils.MapKeysToList(a.config.ParsedAddons), "Done"),
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "Done" {
			break
		}

		err = a.menuAddonSettings(cluster, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *addonClusterMenu) menuAddonSettings(cluster *project.Cluster, addon string) error {
	for {
		selectOptions := []string{"Enable", "Done"}
		if (*cluster).IsAddonEnabled(addon) {
			selectOptions = []string{"Disable", "Settings", "Done"}
		}

		prompt := promptui.Select{
			Label: "Settings",
			Items: selectOptions,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		switch result {
		case "Enable":
			fmt.Println("Enable addon", addon)
			cluster.EnableAddon(addon)
		case "Disable":
			fmt.Println("Disable addon", addon)
			cluster.DisableAddon(addon)
		case "Settings":
			err := a.menuAddonProperties(cluster, addon)
			if err != nil {
				return err
			}
		case "Done":
			return nil
		default:
			return fmt.Errorf("invalid option %s", result)
		}

		// FIXME: remove debug output
		pp.Println(cluster)
	}
}

func (a *addonClusterMenu) menuAddonProperties(cluster *project.Cluster, addon string) error {
	for {
		prompt := promptui.Select{
			Label: "Properties",
			Items: append(utils.MapKeysToList(a.config.ParsedAddons[addon].Properties), "Done"),
			// TODO: add template to display property options
			Templates: a.menuTemplateAddonProperties(cluster, addon),
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "Done" {
			break
		}

		value, err := cli.UntypedQuestion(a.writer, a.reader, "Value", cluster.Addons[addon].Properties[result], func(s any) error {
			if s == nil {
				return fmt.Errorf("value cannot be empty")
			}
			return nil
		})
		if err != nil {
			fmt.Println("Value violates requirements, please try again", err)
			continue
		}
		if cluster.Addons[addon].Properties == nil {
			cluster.Addons[addon].Properties = map[string]any{}
		}
		cluster.Addons[addon].Properties[result] = value
	}
	return nil
}

// menuTemplateAddonProperties returns a promptui.SelectTemplates for the addon properties
func (a *addonClusterMenu) menuTemplateAddonProperties(cluster *project.Cluster, addon string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:   "{{ . }}",
		Details: "{{ properties . }}",
		FuncMap: func() map[string]any {
			funcmap := promptui.FuncMap
			funcmap["properties"] = func(selectValue string) string {
				if selectValue == "Done" {
					return ""
				}
				resultString := "--------------------------------\nDetails:\n"
				resultString += fmt.Sprintf("\tDescription: %s\n", a.config.ParsedAddons[addon].Properties[selectValue].Description)
				resultString += fmt.Sprintf("\tRequired: %v\n", a.config.ParsedAddons[addon].Properties[selectValue].Required)
				resultString += fmt.Sprintf("\tType: %v\n", a.config.ParsedAddons[addon].Properties[selectValue].Type)
				resultString += fmt.Sprintf("\tDefault: %v\n", a.config.ParsedAddons[addon].Properties[selectValue].Default)
				data := cluster.Addons[addon].Properties[selectValue]
				if data == nil {
					data = a.config.ParsedAddons[addon].Properties[selectValue].Default
				}
				resultString += fmt.Sprintf("\tValue: %v\n", data)
				return resultString
			}
			return funcmap
		}(),
	}
}

type addonMenu struct {
	writer io.Writer
	reader *bufio.Reader
	config *project.ProjectConfig
}

func (a *addonMenu) menuCreateAddon() (*project.Addon, error) {
	// ask for the addon name
	addonName, err := cli.StringQuestion(a.writer, a.reader, "Addon Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("addon name cannot be empty")
		}

		if _, ok := a.config.Addons[s]; ok {
			return fmt.Errorf("addon already exists")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// let's ask if the user want to add additional properties
	isEnabledByDefault, err := cli.BooleanQuestion(a.writer, a.reader, "Should this addon be enabled by default?", false)
	if err != nil {
		return nil, err
	}

	sourcePath, err := cli.StringQuestion(a.writer, a.reader, "Please provide the path to the location of the addon (the directory must contain a manifest.yaml file)", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("addon source path cannot be empty")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	prompt := promptui.SelectWithAdd{
		Label:    "Select Group",
		Items:    a.config.AddonGroups(),
		AddLabel: "Create new group",
	}
	_, groupName, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	// let's ask if the user want to add additional properties
	confirmation, err := cli.BooleanQuestion(a.writer, a.reader, "Are you sure you want to create the addon?", false)
	if err != nil {
		return nil, err
	}
	if !confirmation {
		return nil, fmt.Errorf("canceled addon creation")
	}

	return &project.Addon{
		Name:           addonName,
		DefaultEnabled: isEnabledByDefault,
		Path:           sourcePath,
		Group:          groupName,
	}, nil
}

func (a *addonMenu) menuUpdateAddon(addonName string) (*project.Addon, error) {
	return nil, fmt.Errorf("menuUpdateAddon not implemented")
}

func (a *addonMenu) menuDeleteAddon(addonName string) (*project.Addon, error) {
	return nil, fmt.Errorf("menuDeleteAddon not implemented")
}
