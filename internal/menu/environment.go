package menu

import (
	"bufio"
	"fmt"
	"io"
	"path"

	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

// CreateEnvironment creates a context menu to create a new environment
// As part of this we ask for the environment name
func CreateEnvironment(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) (*CarrierCreateEnvironment, error) {
	environmentName, err := cli.StringQuestion(writer, reader, "Environment Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("environment name cannot be empty")
		}
		if _, ok := config.Environments[s]; ok {
			return fmt.Errorf("environment already exists")
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
		pts, err := askForProperties(map[string]string{}, writer, reader)
		if err != nil {
			return nil, err
		}
		properties = pts
	}

	// ask for confirmation
	fqnPath := path.Join(config.BasePath, environmentName)
	confirmation, err := cli.BooleanQuestion(writer, reader, fmt.Sprintf("Are you sure to create a new environment in %s", fqnPath), false)
	if err != nil {
		return nil, err
	}
	if !confirmation {
		return nil, fmt.Errorf("confirmation denied")
	}

	return &CarrierCreateEnvironment{
		EnvironmentName: environmentName,
		Properties:      properties,
	}, nil
}

// UpdateEnvironment creates a context menu to update an existing environment
// As part of this we ask for the environment name
func UpdateEnvironment(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) (*CarrierCreateEnvironment, error) {
	// read environment
	prompt := promptui.Select{
		Label:     "Select Cluster",
		Items:     utils.MapKeysToList(config.Environments),
		Templates: helperSelectTemplate(config, "", ""),
	}
	_, envName, err := prompt.Run()
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
		pts, err := askForProperties(config.Environments[envName].Properties, writer, reader)
		if err != nil {
			return nil, err
		}
		properties = pts
	}

	return &CarrierCreateEnvironment{
		EnvironmentName: envName,
		Properties:      properties,
	}, nil
}
