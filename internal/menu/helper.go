package menu

import (
	"fmt"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/manifoldco/promptui"
)

func helperSelectTemplate(config *project.ProjectConfig, environment, stage string) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "",
		Inactive: "",
		Selected: "",
		Details:  "{{ properties . }}",
		FuncMap: func() map[string]any {
			funcmap := promptui.FuncMap
			funcmap["properties"] = func(selectValue string) string {
				props := map[string]string{}
				addons := map[string]map[string]any{}
				if environment == "" && stage == "" && selectValue != "" {
					props = config.Environments[selectValue].Properties
				}
				if environment != "" && stage == "" && selectValue == "" {
					props = config.EnvStageProperty(environment, selectValue)
				}
				if environment != "" && stage != "" && selectValue != "" {
					addons = config.Environments[environment].Stages[stage].Clusters[selectValue].Addons
					props = config.EnvStageClusterProperty(environment, stage, selectValue)
				}
				resultString := ""
				// add addons to the result string
				if len(addons) > 0 {
					fmt.Println("Addons found")
					resultString += "--------------------------------\nAddons:\n"
					for k, v := range addons {
						resultString += fmt.Sprintf("\t%s:\n", k)
						for kk, vv := range v {
							resultString += fmt.Sprintf("\t\t%s: %v\n", kk, vv)
						}
					}
				}

				// add properties to the result string
				resultString += "--------------------------------\nProperties:\n"
				rslt := []string{}
				for k, v := range props {
					rslt = append(rslt, fmt.Sprintf("%s: %s", k, v))
				}
				for _, v := range rslt {
					resultString += "\t" + v + "\n"
				}
				return resultString
			}
			return funcmap
		}(),
	}
}

func helperSelectTemplateAddonProperties(c project.Cluster, tmpl template.TemplateManifest) *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "",
		Inactive: "",
		Selected: "",
		Details:  "{{ properties . }}",
		FuncMap: func() map[string]any {
			funcmap := promptui.FuncMap
			funcmap["properties"] = func(selectValue string) string {
				resultString := "--------------------------------\nDetails:\n"
				resultString += fmt.Sprintf("\tDescription: %s\n", tmpl.Properties[selectValue].Description)
				resultString += fmt.Sprintf("\tRequired: %v\n", tmpl.Properties[selectValue].Required)
				resultString += fmt.Sprintf("\tType: %v\n", tmpl.Properties[selectValue].Type)
				resultString += fmt.Sprintf("\tDefault: %v\n", tmpl.Properties[selectValue].Default)
				data := c.Addons[tmpl.Name][selectValue]
				if data == nil {
					data = tmpl.Properties[selectValue].Default
				}
				resultString += fmt.Sprintf("\tValue: %v\n", data)
				return resultString
			}
			return funcmap
		}(),
	}
}
