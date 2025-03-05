package menu

import (
	"bufio"
	"fmt"
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
)

type stageMenu struct {
	writer io.Writer
	reader *bufio.Reader
	config *project.ProjectConfig
}

func (s *stageMenu) menuCreateStage(env string) (*project.Stage, error) {
	stage, err := cli.StringQuestion(s.writer, s.reader, "Stage Name", "", func(str string) error {
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

	// TODO: ask for properties

	return &project.Stage{
		Name:       stage,
		Properties: map[string]string{},
		Actions:    project.Actions{},
		Clusters:   map[string]*project.Cluster{},
	}, nil
}

func (s *stageMenu) menuUpdateStage(env, stage string) (*project.Stage, error) {
	// TODO: menu is missing to update properties
	return s.config.Environments[env].Stages[stage], nil
}

func (s *stageMenu) menuDeleteStage(env, stage string) error {
	// TODO: menu is missing to delete the stage (cascade delete)
	return nil
}
