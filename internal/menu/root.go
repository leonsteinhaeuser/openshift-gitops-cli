package menu

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

const (
	rootOptionEnvironment = "Manage Environment"
	rootOptionStage       = "Manage Stage"
	rootOptionCluster     = "Manage Cluster"
	rootOptionAddon       = "Manage Addon"
	rootOptionDone        = "Done"
)

// RootMenu is the main menu of the application
func RootMenu(config *project.ProjectConfig, eventCh chan<- Event) error {
	for {
		prompt := promptui.Select{
			Label: "Action",
			Items: []string{rootOptionEnvironment, rootOptionStage, rootOptionCluster, rootOptionAddon, rootOptionDone},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		switch result {
		case rootOptionEnvironment:
			environmentMenu := &environmentMenu{
				writer: os.Stdout,
				reader: bufio.NewReader(os.Stdin),
				config: config,
			}

			create := func() error {
				environment, err := environmentMenu.menuCreateEnvironment()
				if err != nil {
					return err
				}

				eventCh <- newPreCreateEvent(EventOriginEnvironment, environment.Name, "", "")
				// add environment to config
				config.Environments[environment.Name] = environment
				eventCh <- newPostCreateEvent(EventOriginEnvironment, environment.Name, "", "")
				return nil
			}

			update := func() error {
				env, err := menuSelectEnvironment(config)
				if err != nil {
					return err
				}
				if env == nil {
					return fmt.Errorf("no environment selected")
				}

				environment, err := environmentMenu.menuUpdateEnvironment(*env)
				if err != nil {
					return err
				}

				eventCh <- newPreUpdateEvent(EventOriginEnvironment, environment.Name, "", "")
				config.Environments[*env].Properties = environment.Properties
				eventCh <- newPostUpdateEvent(EventOriginEnvironment, environment.Name, "", "")
				return nil
			}

			delete := func() error {
				env, err := menuSelectEnvironment(config)
				if err != nil {
					return err
				}
				if env == nil {
					return fmt.Errorf("no environment selected")
				}

				_, err = environmentMenu.menuDeleteEnvironment(*env)
				if err != nil {
					return err
				}
				// TODO: add pre and post hooks as well as the actual delete
				return errors.New("environment delete not implemented")
			}

			return crudMenu("Environment Action", create, update, delete)
		case rootOptionStage:
			sm := &stageMenu{
				writer: os.Stdout,
				reader: bufio.NewReader(os.Stdin),
				config: config,
			}

			create := func() error {
				envName, err := menuSelectEnvironment(config)
				if err != nil {
					return err
				}

				stage, err := sm.menuCreateStage(*envName)
				if err != nil {
					return err
				}

				eventCh <- newPreCreateEvent(EventOriginStage, *envName, stage.Name, "")
				// add stage to config
				config.Environments[*envName].Stages[stage.Name] = stage
				eventCh <- newPostCreateEvent(EventOriginStage, *envName, stage.Name, "")
				return nil
			}

			update := func() error {
				envName, stageName, err := menuHierarchySelectEnvironmentStage(config)
				if err != nil {
					return err
				}

				stage, err := sm.menuUpdateStage(*envName, *stageName)
				if err != nil {
					return err
				}

				eventCh <- newPreUpdateEvent(EventOriginStage, *envName, *stageName, "")
				// update stage in config
				config.Environments[*envName].Stages[stage.Name] = stage
				eventCh <- newPostUpdateEvent(EventOriginStage, *envName, stage.Name, "")
				return nil
			}

			delete := func() error {
				envName, stageName, err := menuHierarchySelectEnvironmentStage(config)
				if err != nil {
					return err
				}

				err = sm.menuDeleteStage(*envName, *stageName)
				if err != nil {
					return err
				}

				// TODO: pre and post hooks are missing and the actual delete
				return errors.New("stage delete not implemented")
			}

			return crudMenu("Stage Action", create, update, delete)
		case rootOptionCluster:
			cm := &clusterMenu{
				writer: os.Stdout,
				reader: bufio.NewReader(os.Stdin),
				config: config,
			}

			create := func() error {
				env, stage, err := menuHierarchySelectEnvironmentStage(config)
				if err != nil {
					return err
				}

				cluster, err := cm.menuCreateCluster(*env, *stage)
				if err != nil {
					return err
				}

				eventCh <- newPreCreateEvent(EventOriginCluster, *env, *stage, cluster.Name)
				config.SetCluster(*env, *stage, cluster)
				eventCh <- newPostCreateEvent(EventOriginCluster, *env, *stage, cluster.Name)
				return nil
			}

			update := func() error {
				envName, stageName, clusterName, err := menuHierarchySelectEnvironmentStageCluster(config)
				if err != nil {
					return err
				}

				cluster, err := cm.menuUpdateCluster(*envName, *stageName, *clusterName)
				if err != nil {
					return err
				}

				eventCh <- newPreUpdateEvent(EventOriginCluster, *envName, *stageName, *clusterName)
				config.SetCluster(*envName, *stageName, cluster)
				eventCh <- newPostUpdateEvent(EventOriginCluster, *envName, *stageName, *clusterName)
				return nil
			}

			delete := func() error {
				envName, stageName, clusterName, err := menuHierarchySelectEnvironmentStageCluster(config)
				if err != nil {
					return err
				}

				_, err = cm.menuDeleteCluster(*envName, *stageName, *clusterName)
				if err != nil {
					return err
				}

				eventCh <- newPreDeleteEvent(EventOriginCluster, *envName, *stageName, *clusterName)
				config.DeleteCluster(*envName, *stageName, *clusterName)
				eventCh <- newPostDeleteEvent(EventOriginCluster, *envName, *stageName, *clusterName)
				return nil
			}
			return crudMenu("Cluster Action", create, update, delete)
		case rootOptionAddon:
			addonMenu := addonMenu{
				writer: os.Stdout,
				reader: bufio.NewReader(os.Stdin),
				config: config,
			}

			create := func() error {
				addon, err := addonMenu.menuCreateAddon()
				if err != nil {
					return err
				}

				eventCh <- newPreCreateEvent(EventOriginAddon, addon.Name, "", "")
				// add addon to config
				config.Addons[addon.Name] = *addon
				eventCh <- newPostCreateEvent(EventOriginAddon, addon.Name, "", "")
				return nil
			}

			update := func() error {
				addonName, err := menuSelectAddon(config)
				if err != nil {
					return err
				}

				eventCh <- newPreUpdateEvent(EventOriginAddon, *addonName, "", "")
				addon, err := addonMenu.menuUpdateAddon(*addonName)
				if err != nil {
					return err
				}
				config.Addons[*addonName] = *addon
				eventCh <- newPostUpdateEvent(EventOriginAddon, *addonName, "", "")
				return nil
			}

			delete := func() error {
				addonName, err := menuSelectAddon(config)
				if err != nil {
					return err
				}

				_, err = addonMenu.menuDeleteAddon(*addonName)
				if err != nil {
					return err
				}
				return errors.New("addon delete not implemented")
			}

			return crudMenu("Addon Action", create, update, delete)
		case rootOptionDone:
			return nil
		default:
			return fmt.Errorf("invalid option %s", result)
		}
	}
}

// menuSelectAddon is a helper function to select an addon
func menuSelectAddon(config *project.ProjectConfig) (*string, error) {
	prompt := promptui.Select{
		Label: "Select Addon",
		Items: append(utils.MapKeysToList(config.Addons), rootOptionDone),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &result, err
}

// menuSelectEnvironment is a helper function to select an environment
func menuSelectEnvironment(config *project.ProjectConfig) (*string, error) {
	prompt := promptui.Select{
		Label: "Select Environment",
		Items: append(utils.MapKeysToList(config.Environments), rootOptionDone),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &result, err
}

// menuSelectStage is a helper function to select a stage
func menuSelectStage(config *project.ProjectConfig, environment string) (*string, error) {
	prompt := promptui.Select{
		Label: "Select Stage",
		Items: append(utils.MapKeysToList(config.Environments[environment].Stages), rootOptionDone),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &result, err
}

// menuHierarchySelectEnvironmentStage is a helper function to select an environment and stage
func menuHierarchySelectEnvironmentStage(config *project.ProjectConfig) (*string, *string, error) {
	env, err := menuSelectEnvironment(config)
	if err != nil {
		return nil, nil, err
	}
	if env == nil {
		return nil, nil, nil
	}

	stage, err := menuSelectStage(config, *env)
	if err != nil {
		return nil, nil, err
	}
	if env == nil {
		return nil, nil, nil
	}

	return env, stage, nil
}

// menuSelectCluster is a helper function to select a cluster
func menuSelectCluster(config *project.ProjectConfig, environment, stage string) (*string, error) {
	prompt := promptui.Select{
		Label: "Select Cluster",
		Items: append(utils.MapKeysToList(config.Environments[environment].Stages[stage].Clusters), rootOptionDone),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &result, err
}

// menuHierarchySelectEnvironmentStage is a helper function to select an environment and stage
func menuHierarchySelectEnvironmentStageCluster(config *project.ProjectConfig) (*string, *string, *string, error) {
	env, stage, err := menuHierarchySelectEnvironmentStage(config)
	if err != nil {
		return nil, nil, nil, err
	}

	cluster, err := menuSelectCluster(config, *env, *stage)
	if err != nil {
		return nil, nil, nil, err
	}
	if env == nil {
		return nil, nil, nil, nil
	}
	return env, stage, cluster, nil
}
