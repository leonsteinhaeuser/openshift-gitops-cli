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

type environmentMenu struct {
	writer io.Writer
	reader *bufio.Reader
	config *project.ProjectConfig
}

func (e *environmentMenu) menuCreateEnvironment() (*project.Environment, error) {
	env, err := cli.StringQuestion(e.writer, e.reader, "Environment Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("environment name cannot be empty")
		}
		if _, ok := e.config.Environments[s]; ok {
			return fmt.Errorf("environment already exists")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	environment := &project.Environment{
		Name:       env,
		Stages:     map[string]*project.Stage{},
		Properties: map[string]string{},
	}

	// ask for properties
	properties, err := e.menuEnvironmentProperties(environment)
	if err != nil {
		return nil, err
	}
	environment.Properties = properties

	return environment, nil
}

func (e *environmentMenu) menuUpdateEnvironment(envName string) (*project.Environment, error) {
	environment := *e.config.Environments[envName]
	environment.Name = envName
	// ask for properties
	properties, err := e.menuEnvironmentProperties(&environment)
	if err != nil {
		return nil, err
	}
	environment.Properties = properties
	return &environment, nil
}

func (e *environmentMenu) menuDeleteEnvironment(envName string) (*project.Environment, error) {
	confirmation, err := cli.BooleanQuestion(e.writer, e.reader, fmt.Sprintf("Are you sure to delete the environment %s. Keep in mind that this will also delete the stages and clusters.", envName), false)
	if err != nil {
		return nil, err
	}
	if !confirmation {
		return nil, fmt.Errorf("confirmation denied")
	}
	environment := *e.config.Environments[envName]
	environment.Name = envName
	return &environment, errors.New("menuDeleteEnvironment not implemented")
}

func (e *environmentMenu) menuEnvironmentProperties(env *project.Environment) (map[string]string, error) {
	envProperties := map[string]string{}
	for {
		properties := utils.MergeMaps(env.Properties, envProperties)

		prompt := promptui.SelectWithAdd{
			Label:    "Properties",
			Items:    append(utils.MapKeysToList(properties), "Done"),
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

		val, err := cli.StringQuestion(e.writer, e.reader, "Property Value", properties[result], func(s string) error {
			if s == "" {
				return fmt.Errorf("property value cannot be empty")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		envProperties[result] = val
	}
	return envProperties, nil
}
