package menu

import (
	"bufio"
	"fmt"
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
)

func AddAddon(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) error {
	// ask for the addon name
	addonName, err := cli.StringQuestion(writer, reader, "Addon Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("addon name cannot be empty")
		}

		if _, ok := config.Addons[s]; ok {
			return fmt.Errorf("addon already exists")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// let's ask if the user want to add additional properties
	isEnabledByDefault, err := cli.BooleanQuestion(writer, reader, "Should this addon be enabled by default?", false)
	if err != nil {
		return err
	}

	sourcePath, err := cli.StringQuestion(writer, reader, "Please provide the path to the location of the addon (the directory must contain a manifest.yaml file)", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("addon source path cannot be empty")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// let's ask if the user want to add additional properties
	confirmation, err := cli.BooleanQuestion(writer, reader, "Are you sure you want to create the addon?", false)
	if err != nil {
		return err
	}
	if !confirmation {
		return fmt.Errorf("canceled addon creation")
	}

	if config.Addons == nil {
		config.Addons = map[string]project.Addon{}
	}

	config.Addons[addonName] = project.Addon{
		DefaultEnabled: isEnabledByDefault,
		Path:           sourcePath,
	}
	return nil
}
