package project

import (
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/utils"
)

type ProjectConfig struct {
	BasePath         string                 `json:"basePath"`
	TemplateBasePath string                 `json:"templateBasePath"`
	Environments     map[string]Environment `json:"environments"`
}

type Environment struct {
	Name       string            `json:"-"`
	Properties map[string]string `json:"properties"`
	Actions    Actions           `json:"actions"`
	Stages     map[string]Stage  `json:"stages"`
}

type Stage struct {
	Name       string             `json:"-"`
	Properties map[string]string  `json:"properties"`
	Actions    Actions            `json:"actions"`
	Clusters   map[string]Cluster `json:"clusters"`
}

type Cluster struct {
	Name       string            `json:"-"`
	Properties map[string]string `json:"properties"`
}

// EnvStageProperty merges the properties of the environment and stage and returns them as a map
func (pc *ProjectConfig) EnvStageProperty(environment, stage string) map[string]string {
	return utils.MergeMaps(pc.Environments[environment].Properties, pc.Environments[environment].Stages[stage].Properties)
}

func (pc *ProjectConfig) EnvStageClusterProperty(environment, stage, cluster string) map[string]string {
	return utils.MergeMaps(pc.EnvStageProperty(environment, stage), pc.Environments[environment].Stages[stage].Clusters[cluster].Properties)
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

type Command struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// execute executes the command with the given arguments
func (c Command) execute(stdout, errout io.Writer) error {
	return utils.ExecuteShellCommand(stdout, errout, c.Command, c.Args...)
}
