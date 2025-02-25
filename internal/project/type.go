package project

import "github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/utils"

type ProjectConfig struct {
	BasePath         string                 `json:"basePath"`
	TemplateBasePath string                 `json:"templateBasePath"`
	Environments     map[string]Environment `json:"environments"`
}

type Environment struct {
	Name       string            `json:"-"`
	Properties map[string]string `json:"properties"`
	Stages     map[string]Stage  `json:"stages"`
}

type Stage struct {
	Name       string             `json:"-"`
	Properties map[string]string  `json:"properties"`
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
