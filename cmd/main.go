package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/menu"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
)

var (
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
			fmt.Printf("An error occurred while loading the addon [%s] manifest file: %s, %v\n", k, v.Path, err)
			os.Exit(1)
			return
		}
		tm.Name = k
		tm.BasePath = v.Path
		tm.Group = v.Group
		projectConfig.ParsedAddons[k] = *tm
	}
}

func main() {
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
					continue
				}

				if event.Environment != "" && event.Stage == "" && event.Cluster == "" {
					env := projectConfig.Environments[event.Environment]
					err := executeHook(os.Stdout, os.Stderr, event.Type, event.Runtime, env.Actions)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if event.Environment != "" && event.Stage != "" && event.Cluster == "" {
					stage := projectConfig.Environments[event.Environment].Stages[event.Stage]
					err := executeHook(os.Stdout, os.Stderr, event.Type, event.Runtime, stage.Actions)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if event.Environment != "" && event.Stage != "" && event.Cluster != "" {
					cluster := projectConfig.Cluster(event.Environment, event.Stage, event.Cluster)
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
		fmt.Println(err)
		os.Exit(1)
	}
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
