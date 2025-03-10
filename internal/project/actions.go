package project

import (
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
)

type Command struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// execute executes the command with the given arguments
func (c Command) execute(stdout, errout io.Writer) error {
	return utils.ExecuteShellCommand(stdout, errout, c.Command, c.Args...)
}

type Actions struct {
	PreCreateHooks  []Command `json:"preCreateHooks"`
	PostCreateHooks []Command `json:"postCreateHooks"`
	PreUpdateHooks  []Command `json:"preUpdateHooks"`
	PostUpdateHooks []Command `json:"postUpdateHooks"`
}

func executeCommands(stdout, errout io.Writer, commands []Command) error {
	for _, c := range commands {
		err := c.execute(stdout, errout)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a Actions) ExecutePreCreateHooks(stdout, errout io.Writer) error {
	return executeCommands(stdout, errout, a.PreCreateHooks)
}

func (a Actions) ExecutePostCreateHooks(stdout, errout io.Writer) error {
	return executeCommands(stdout, errout, a.PostCreateHooks)
}

func (a Actions) ExecutePreUpdateHooks(stdout, errout io.Writer) error {
	return executeCommands(stdout, errout, a.PreUpdateHooks)
}

func (a Actions) ExecutePostUpdateHooks(stdout, errout io.Writer) error {
	return executeCommands(stdout, errout, a.PostUpdateHooks)
}
