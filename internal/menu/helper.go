package menu

import (
	"fmt"

	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/project"
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
				if environment == "" && stage == "" && selectValue != "" {
					props = config.Environments[selectValue].Properties
				}
				if environment != "" && stage == "" && selectValue == "" {
					props = config.EnvStageProperty(environment, selectValue)
				}
				if environment != "" && stage != "" && selectValue != "" {
					props = config.EnvStageClusterProperty(environment, stage, selectValue)
				}

				resultString := "--------------------------------\nProperties:\n"
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
