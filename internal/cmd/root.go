package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/menu"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/spf13/cobra"
)

const (
	flagKeyEnvironment = "environment"
	flagKeyStage       = "stage"
	flagKeyCluster     = "cluster"
	flagKeyAddon       = "addon"

	PROJECTFILENAME = "PROJECT.yaml"
)

var (
	RootCmd = &cobra.Command{
		Use: "root",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd()
		},
	}

	projectConfig *project.ProjectConfig
)

func init() {
	// load project configuration
	pc, err := project.ParseConfig(PROJECTFILENAME)
	if err != nil {
		fmt.Println("An error occurred while parsing the PROJECT.yaml file", err)
		os.Exit(1)
		return
	}
	projectConfig = pc

	// initialize the run command
	run := &cobra.Command{
		Use:   "run",
		Short: "Render the cluster configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCmd(cmd)
		},
	}
	run.PersistentFlags().StringP(flagKeyEnvironment, "e", "", "environment the check should be performed for")
	run.PersistentFlags().StringP(flagKeyStage, "s", "", "stage the check should be performed for")
	run.PersistentFlags().StringP(flagKeyCluster, "c", "", "cluster the check should be performed for")
	RootCmd.AddCommand(run)
}

func rootCmd() error {
	eventsPipeline := make(chan menu.Event, 100)
	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	go func(ctx context.Context) {
		// handle config file updates
		for {
			select {
			case <-ctx.Done():
				close(eventsPipeline)
				return
			case event := <-eventsPipeline:
				// we only need to update the config file if the action is a post action
				// because we need to update the config only, if the action was successful
				if event.Runtime == menu.EventRuntimePost {
					// update config file
					err := project.UpdateOrCreateConfig(PROJECTFILENAME, projectConfig)
					if err != nil {
						fmt.Println("An error occurred while updating the project config", err)
						return
					}
				}

				if event.Origin == menu.EventOriginAddon {
					if event.Runtime == menu.EventRuntimePre {
						addonPath := projectConfig.Addons[event.Environment].Path
						_, err := template.LoadManifest(projectConfig.Addons[event.Environment].Path)
						if err != nil {
							fmt.Printf("An error occurred while loading the addon [%s] manifest file: %s, %v\n", event.Environment, addonPath, err)
							os.Exit(1)
							return
						}
					}
					continue
				}

				if event.Environment != "" && event.Stage == "" && event.Cluster == "" {
					env := projectConfig.GetEnvironment(event.Environment)
					err := executeHook(os.Stdout, os.Stderr, event.Type, event.Runtime, env.Actions)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if event.Environment != "" && event.Stage != "" && event.Cluster == "" {
					stage := projectConfig.GetStage(event.Environment, event.Stage)
					err := executeHook(os.Stdout, os.Stderr, event.Type, event.Runtime, stage.Actions)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if event.Environment != "" && event.Stage != "" && event.Cluster != "" {
					cluster := projectConfig.GetCluster(event.Environment, event.Stage, event.Cluster)
					if event.Type == menu.EventTypeCreate || event.Type == menu.EventTypeUpdate {
						err := cluster.Render(projectConfig, event.Environment, event.Stage)
						if err != nil {
							fmt.Printf("An error occurred while rendering the cluster [%s] configuration: %v", event.Cluster, err)
							return
						}
					}
				}
			}
		}
	}(ctx)

	err := menu.RootMenu(projectConfig, eventsPipeline)
	if err != nil {
		return err
	}
	return nil
}

func executeHook(stdout, errout io.Writer, t menu.EventType, r menu.EventRuntime, actions project.Actions) error {
	switch t {
	case menu.EventTypeCreate:
		if r == menu.EventRuntimePre {
			err := actions.ExecutePreCreateHooks(stdout, errout)
			if err != nil {
				return err
			}
		}
		if r == menu.EventRuntimePost {
			err := actions.ExecutePostCreateHooks(stdout, errout)
			if err != nil {
				return err
			}
		}
		return nil
	case menu.EventTypeUpdate:
		if r == menu.EventRuntimePre {
			err := actions.ExecutePreUpdateHooks(stdout, errout)
			if err != nil {
				return err
			}
		}
		if r == menu.EventRuntimePost {
			err := actions.ExecutePostUpdateHooks(stdout, errout)
			if err != nil {
				return err
			}
		}
	case menu.EventTypeDelete:
	default:
		return fmt.Errorf("unknown event type: %v", t)
	}
	return nil
}

// runCmd expects the environment, stage and cluster flags to be set and renders the cluster with all properties and addons to disk
func runCmd(cmd *cobra.Command) error {
	env, err := cmd.Flags().GetString(flagKeyEnvironment)
	if err != nil {
		return fmt.Errorf("flag environment not found %w", err)
	}
	stg, err := cmd.Flags().GetString(flagKeyStage)
	if err != nil {
		return fmt.Errorf("flag stage not found: %w", err)
	}
	cluster, err := cmd.Flags().GetString(flagKeyCluster)
	if err != nil {
		return fmt.Errorf("flag cluster not found %w", err)
	}

	switch {
	case env == "":
		return fmt.Errorf("environment is required")
	case stg == "":
		return fmt.Errorf("stage is required")
	case cluster == "":
		return fmt.Errorf("cluster is required")
	}

	clusterCfg := projectConfig.GetCluster(env, stg, cluster)

	err = clusterCfg.Render(projectConfig, env, stg)
	if err != nil {
		return err
	}
	return nil
}
