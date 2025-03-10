package menu

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

type stageMenu struct {
	writer io.Writer
	reader *bufio.Reader
	config *project.ProjectConfig
}

func (s *stageMenu) menuCreateStage(env string) (*project.Stage, error) {
	stageName, err := cli.StringQuestion(s.writer, s.reader, "Stage Name", "", func(str string) error {
		if str == "" {
			return fmt.Errorf("stage name cannot be empty")
		}
		if _, ok := s.config.Environments[env].Stages[str]; ok {
			return fmt.Errorf("stage already exists")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	stage := &project.Stage{
		Name:       stageName,
		Properties: map[string]string{},
		Actions:    project.Actions{},
		Clusters:   map[string]*project.Cluster{},
	}

	properties, err := s.menuProperties(stage)
	if err != nil {
		return nil, err
	}
	stage.Properties = properties

	return stage, nil
}

func (s *stageMenu) menuUpdateStage(envName, stageName string) (*project.Stage, error) {
	stage := s.config.Environments[envName].Stages[stageName]
	properties, err := s.menuProperties(stage)
	if err != nil {
		return nil, err
	}
	stage.Properties = properties
	return stage, nil
}

func (s *stageMenu) menuDeleteStage(env, stage string) error {
	// TODO: menu is missing to delete the stage (cascade delete)
	return errors.New("not implemented")
}

func (s *stageMenu) menuProperties(stage *project.Stage) (map[string]string, error) {
	stageProperties := map[string]string{}
	for {
		properties := utils.MergeMaps(stage.Properties, stageProperties)

		prompt := promptui.SelectWithAdd{
			Label:    "Properties",
			Items:    append(utils.SortStringSlice(utils.MapKeysToList(properties)), "Done"),
			AddLabel: "Create Property",
		}
		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		if result == "" {
			return nil, fmt.Errorf("property key cannot be empty")
		}
		if result == "Done" {
			// user is done
			break
		}

		val, err := cli.StringQuestion(s.writer, s.reader, "Property Value", properties[result], func(s string) error {
			if s == "" {
				return fmt.Errorf("property value cannot be empty")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		stageProperties[result] = val
	}
	return stageProperties, nil
}
