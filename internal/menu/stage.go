package menu

import (
	"bufio"
	"fmt"
	"io"
	"path"

	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

// CreateStage creates a context menu to create a new stage
// As part of this we ask for the environment and stage name
func CreateStage(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) (*CarrierCreateStage, error) {
	// read environments
	prompt := promptui.Select{
		Label: "Select Environment",
		Items: utils.MapKeysToList(config.Environments),
	}
	_, envResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	stageName, err := cli.StringQuestion(writer, reader, "Stage Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("stage name cannot be empty")
		}
		if _, ok := config.Environments[envResult].Stages[s]; ok {
			return fmt.Errorf("stage already exists")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// ask for confirmation
	fqnPath := path.Join(config.BasePath, envResult, stageName)
	confirmation, err := cli.BooleanQuestion(writer, reader, fmt.Sprintf("Are you sure to create a new stage in %s", fqnPath), false)
	if err != nil {
		return nil, err
	}
	if !confirmation {
		return nil, fmt.Errorf("confirmation denied")
	}
	return &CarrierCreateStage{
		Environment: envResult,
		StageName:   stageName,
	}, nil
}
