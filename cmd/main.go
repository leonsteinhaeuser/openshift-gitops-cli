package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/menu"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

var (
	actionMap = map[string]func() error{
		"Create Cluster": func() error {
			ccc, err := menu.CreateCluster(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			err = projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Actions.ExecutePreCreateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute pre create hooks: %w", err)
			}

			cprops := ccc.Properties
			// merge properties from environment and stage
			props := utils.MergeMaps(projectConfig.Environments[ccc.Environment].Properties, projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Properties, ccc.Properties)
			ccc.Properties = props

			projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Clusters[ccc.ClusterName] = project.Cluster{
				Properties: cprops,
				Addons:     ccc.Addons,
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
				err = t.Render(projectConfig.BasePath, template.TemplateData{
					Environment: ccc.Environment,
					Stage:       ccc.Stage,
					ClusterName: ccc.ClusterName,
					Properties:  ccc.Properties,
				})
				if err != nil {
					return err
				}
			}

			err = projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Actions.ExecutePostCreateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute post create hooks: %w", err)
			}
			return nil
		},
		"Update Cluster": func() error {
			ccc, err := menu.UpdateCluster(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			err = projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Actions.ExecutePreUpdateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute pre update hooks: %w", err)
			}

			cprops := ccc.Properties
			// merge properties from environment and stage
			props := utils.MergeMaps(projectConfig.Environments[ccc.Environment].Properties, projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Properties, ccc.Properties)
			ccc.Properties = props

			projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Clusters[ccc.ClusterName] = project.Cluster{
				Properties: cprops,
			}
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}

			err = projectConfig.Environments[ccc.Environment].Stages[ccc.Stage].Actions.ExecutePostUpdateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute post update hooks: %w", err)
			}

			// TODO: we might need to update the templates
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

			projectConfig.Environments[env.EnvironmentName] = project.Environment{
				Properties: env.Properties,
				Stages:     map[string]project.Stage{},
			}
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}
			return nil
		},
		"Update Environment": func() error {
			env, err := menu.UpdateEnvironment(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			envC := projectConfig.Environments[env.EnvironmentName]
			envC.Properties = env.Properties
			projectConfig.Environments[env.EnvironmentName] = envC
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}
			return nil
		},
		"Create Stage": func() error {
			cc, err := menu.CreateStage(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			err = projectConfig.Environments[cc.Environment].Actions.ExecutePreCreateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute pre create hooks: %w", err)
			}

			if projectConfig.Environments[cc.Environment].Stages == nil {
				projectConfig.Environments[cc.Environment] = project.Environment{
					Stages: map[string]project.Stage{},
				}
			}

			envC := projectConfig.Environments[cc.Environment].Stages[cc.StageName]
			envC.Properties = cc.Properties
			projectConfig.Environments[cc.Environment].Stages[cc.StageName] = envC
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}

			err = projectConfig.Environments[cc.Environment].Actions.ExecutePostCreateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute post create hooks: %w", err)
			}
			return nil
		},
		"Update Stage": func() error {
			cc, err := menu.UpdateStage(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			err = projectConfig.Environments[cc.Environment].Actions.ExecutePreUpdateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute pre update hooks: %w", err)
			}

			projectConfig.Environments[cc.Environment].Stages[cc.StageName] = project.Stage{
				Properties: cc.Properties,
				Clusters:   map[string]project.Cluster{},
			}
			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}

			err = projectConfig.Environments[cc.Environment].Actions.ExecutePostUpdateHooks(os.Stdout, os.Stderr)
			if err != nil {
				return fmt.Errorf("failed to execute post update hooks: %w", err)
			}
			return nil
		},
		"Add Addon": func() error {
			err := menu.AddAddon(projectConfig, os.Stdout, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			err = project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
			if err != nil {
				return err
			}
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
		defer f.Close()
		_, err = f.WriteString("basePath: overlays/\ntemplateBasePath: templates/\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
	}

	pc, err := project.ParseConfig(PROJECTFILENAME)
	if err != nil {
		fmt.Println("An error occurred while parsing the PROJECT.yaml file", err)
		os.Exit(1)
		return
	}
	projectConfig = pc

	if projectConfig.ParsedAddons == nil {
		projectConfig.ParsedAddons = map[string]template.TemplateManifest{}
	}

	// load all addons, so we can use them later
	for k, v := range projectConfig.Addons {
		tm, err := template.LoadManifest(v.Path)
		if err != nil {
			fmt.Printf("An error occurred while loading the addon [%s] manifest file: %s, %v", k, v.Path, err)
			return
		}
		tm.Name = k
		tm.BasePath = v.Path
		projectConfig.ParsedAddons[k] = *tm
	}
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
