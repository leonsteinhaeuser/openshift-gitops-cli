package menu

import (
	"bufio"
	"fmt"
	"io"

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

func (a *addonClusterMenu) menuManageAddons(ah project.AddonHandler) error {
	for {
		prompt := promptui.Select{
			Label:     "Manage Addons",
			Items:     append(utils.SortStringSlice(utils.MapKeysToList(a.config.ParsedAddons)), "Done"),
			Templates: a.templateManageAddons(ah),
			Size:      10,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "Done" {
			err := ah.GetAddons().AllRequiredPropertiesSet(a.config)
			if err != nil {
				fmt.Println("Not all required properties are set", err)
				continue
			}
			fmt.Println("Done managing addons")
			break
		}

		err = a.menuAddonSettings(ah, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *addonClusterMenu) templateManageAddons(ah project.AddonHandler) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:   "{{ . }}",
		Details: "{{ addon . }}",
		FuncMap: func() map[string]any {
			funcmap := promptui.FuncMap
			funcmap["addon"] = func(addonName string) string {
				if addonName == "Done" {
					return ""
				}
				resultString := "--------------------------------\nDetails:\n"
				resultString += fmt.Sprintf("\tDescription: %s\n", a.config.ParsedAddons[addonName].Description)
				resultString += fmt.Sprintf("\tGroup: %s\n", a.config.ParsedAddons[addonName].Group)
				resultString += fmt.Sprintf("\tEnabled: %v | Default: %v\n", ah.GetAddon(addonName).IsEnabled(), a.config.Addons[addonName].DefaultEnabled)
				return resultString
			}
			return funcmap
		}(),
	}
}

func (a *addonClusterMenu) menuAddonSettings(ah project.AddonHandler, addon string) error {
	for {
		selectOptions := []string{"Enable", "Done"}
		if ah.IsAddonEnabled(addon) {
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
			ah.EnableAddon(addon)
		case "Disable":
			ah.DisableAddon(addon)
		case "Settings":
			err := a.menuAddonProperties(ah, addon)
			if err != nil {
				return err
			}
		case "Done":
			if !ah.IsAddonEnabled(addon) {
				return nil
			}

			// check if all required properties are set
			err := ah.GetAddon(addon).AllRequiredPropertiesSet(a.config, addon)
			if err != nil {
				fmt.Println(utils.Red.Wrap("Not all required properties are set"), err)
				continue
			}
			return nil
		default:
			return fmt.Errorf("invalid option %s", result)
		}
	}
}

func (a *addonClusterMenu) menuAddonProperties(ah project.AddonHandler, addon string) error {
	for {
		prompt := promptui.Select{
			Label: "Properties",
			Items: append(utils.SortStringSlice(utils.MapKeysToList(a.config.ParsedAddons[addon].Properties)), "Done"),
			// TODO: add template to display property options
			Templates: a.menuTemplateAddonProperties(ah, addon),
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "Done" {
			break
		}

		value, err := cli.UntypedQuestion(a.writer, a.reader, "Value", fmt.Sprintf("%v", ah.GetAddon(addon).Properties[result]), func(s any) error {
			if s == nil {
				return fmt.Errorf("value cannot be empty")
			}
			return nil
		})
		if err != nil {
			fmt.Println(utils.Red.Wrap("Value violates requirements, please try again"), err)
			continue
		}

		value, err = a.config.ParsedAddons[addon].Properties[result].ParseValue(value)
		if err != nil {
			fmt.Println(utils.Red.Wrap("Value violates requirements, please try again"), err)
			continue
		}
		ah.GetAddon(addon).SetProperty(result, value)
	}
	return nil
}

// menuTemplateAddonProperties returns a promptui.SelectTemplates for the addon properties
func (a *addonClusterMenu) menuTemplateAddonProperties(ah project.AddonHandler, addon string) *promptui.SelectTemplates {
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
				data := ah.GetAddon(addon).Properties[selectValue]
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
	for {
		prompt := promptui.SelectWithAdd{
			Label:    "Select Option",
			Items:    []string{"Set Group", "Set Path", "Enable / Disable by Default", "Done"},
			AddLabel: "Create new group",
		}
		_, option, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		addon := a.config.Addons[addonName]
		switch option {
		case "Set Group":
			prompt := promptui.SelectWithAdd{
				Label:    "Select Group",
				Items:    a.config.AddonGroups(),
				AddLabel: "Create new group",
			}
			_, groupName, err := prompt.Run()
			if err != nil {
				return nil, err
			}
			addon.Group = groupName
			continue
		case "Set Path":
			sourcePath, err := cli.StringQuestion(a.writer, a.reader, "Please provide the path to the location of the addon (the directory must contain a manifest.yaml file)", addon.Path, func(s string) error {
				if s == "" {
					return fmt.Errorf("addon source path cannot be empty")
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			addon.Path = sourcePath
			continue
		case "Enable / Disable by Default":
			isEnabledByDefault, err := cli.BooleanQuestion(a.writer, a.reader, "Should this addon be enabled by default?", addon.DefaultEnabled)
			if err != nil {
				return nil, err
			}
			addon.DefaultEnabled = isEnabledByDefault
			continue
		case "Done":
			return &addon, nil
		default:
			return nil, fmt.Errorf("invalid option %s", option)
		}
	}
}

func (a *addonMenu) menuDeleteAddon(addonName string) (*project.Addon, error) {
	return nil, fmt.Errorf("menuDeleteAddon not implemented")
}
