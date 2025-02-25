package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/menu"
	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/template"
	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

var (
	actionMap = map[string]func() error{
		"Create Cluster": func() error {
			ccc, err := menu.CreateCluster(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			// TODO: add the cluster to the project config
			projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Clusters[ccc.ClusterName] = project.Cluster{
				Properties: ccc.Properties,
			}
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}

			templates, err := template.LoadTemplateManifest(projectConfig.TemplateBasePath)
			if err != nil {
				return err
			}

			for _, t := range templates {
				err = t.Render(projectConfig.BasePath, *ccc)
				if err != nil {
					return err
				}
			}
			return nil
		},
		"Create Environment": func() error {
			env, err := menu.CreateEnvironment(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if projectConfig.Environments == nil {
				projectConfig.Environments = map[string]project.Environment{}
			}

			projectConfig.Environments[*env] = project.Environment{
				Stages: map[string]project.Stage{},
			}
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}
			// TODO: use the returned carrier to create the environment
			return nil
		},
		"Create Stage": func() error {
			cc, err := menu.CreateStage(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if projectConfig.Environments[cc.Environment].Stages == nil {
				projectConfig.Environments[cc.Environment] = project.Environment{
					Stages: map[string]project.Stage{},
				}
			}

			projectConfig.Environments[cc.Environment].Stages[cc.StageName] = project.Stage{
				Clusters: map[string]project.Cluster{},
			}
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}
			// TODO: use the returned carrier to create the stage
			return nil
		},
	}

	projectConfig = &project.ProjectConfig{}
)

const (
	PROJECTFILENAME = "PROJECT.yaml"
)

// check for project file and load it
func init() {
	_, err := os.Stat(PROJECTFILENAME)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println("An error occurred while checking for the PROJECT.yaml file", err)
		os.Exit(1)
		return
	}
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(PROJECTFILENAME)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		f.Close()
	}

	pc, err := project.ParseConfig(PROJECTFILENAME)
	if err != nil {
		fmt.Println("An error occurred while parsing the PROJECT.yaml file", err)
		os.Exit(1)
		return
	}
	projectConfig = pc
}

func main() {
	s := utils.MapKeysToList(actionMap)
	slices.Sort(s)

	prompt := promptui.Select{
		Label: "Select Action",
		Items: s,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	err = actionMap[result]()
	if err != nil {
		fmt.Printf("Action failed %v\n", err)
		return
	}
}
