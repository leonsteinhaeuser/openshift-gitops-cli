package menu

import (
	"bufio"
	"fmt"
	"io"
	"path"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
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

	// let's ask if the user want to add additional properties
	createProperties, err := cli.BooleanQuestion(writer, reader, "Do you want to add properties?", false)
	if err != nil {
		return nil, err
	}
	properties := map[string]string{}
	if createProperties {
		pts, err := askForProperties(config.Environments[envResult].Properties, writer, reader)
		if err != nil {
			return nil, err
		}
		properties = pts
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
		Properties:  utils.ReduceMap(properties, config.Environments[envResult].Properties),
	}, nil
}

// UpdateStage creates a context menu to update an existing stage
// As part of this we ask for the environment name and stage name
func UpdateStage(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) (*CarrierCreateStage, error) {
	// read environments
	prompt := promptui.Select{
		Label: "Select Environment",
		Items: utils.MapKeysToList(config.Environments),
	}
	_, envResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	// read stage
	prompt = promptui.Select{
		Label: "Select Stage",
		Items: utils.MapKeysToList(config.Environments[envResult].Stages),
	}
	_, stageResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	// let's ask if the user want to add additional properties
	createProperties, err := cli.BooleanQuestion(writer, reader, "Do you want to update properties?", false)
	if err != nil {
		return nil, err
	}
	properties := map[string]string{}
	if createProperties {
		pts, err := askForProperties(config.Environments[envResult].Stages[stageResult].Properties, writer, reader)
		if err != nil {
			return nil, err
		}
		properties = pts
	}

	return &CarrierCreateStage{
		Environment: envResult,
		StageName:   stageResult,
		Properties:  utils.ReduceMap(properties, config.Environments[envResult].Properties),
	}, nil
}
