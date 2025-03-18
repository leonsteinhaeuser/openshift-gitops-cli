package project

import (
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
)

type Addon struct {
	Name           string `json:"-"`
	Group          string `json:"group"`
	DefaultEnabled bool   `json:"defaultEnabled"`
	Path           string `json:"path"`
}

type ProjectConfig struct {
	BasePath         string                               `json:"basePath"`
	TemplateBasePath string                               `json:"templateBasePath"`
	Addons           map[string]Addon                     `json:"addons"`
	ParsedAddons     map[string]template.TemplateManifest `json:"-"`
	Environments     map[string]*Environment              `json:"environments"`
}

// HasCluster checks if a cluster exists in the given environment and stage
func (p ProjectConfig) HasCluster(env, stage, cluster string) bool {
	_, ok := p.GetStage(env, stage).Clusters[cluster]
	return ok
}

// SetCluster sets the cluster for the given environment and stage
func (p *ProjectConfig) SetCluster(env, stage string, cluster *Cluster) {
	if p.GetStage(env, stage).Clusters == nil {
		p.GetStage(env, stage).Clusters = map[string]*Cluster{}
	}
	p.GetStage(env, stage).Clusters[cluster.Name] = cluster
}

func (p *ProjectConfig) DeleteCluster(env, stage, cluster string) {
	delete(p.GetStage(env, stage).Clusters, cluster)
}

// EnvStageProperty merges the properties of the environment and stage and returns them as a map
func (pc *ProjectConfig) EnvStageProperty(environment, stage string) map[string]string {
	return utils.MergeMaps(pc.GetEnvironment(environment).Properties, pc.GetStage(environment, stage).Properties)
}

// AddonGroups returns a list of addon groups that have been defined in the addons
func (p ProjectConfig) AddonGroups() []string {
	groups := map[string]bool{}
	for _, a := range p.Addons {
		if a.Group == "" {
			continue
		}
		groups[a.Group] = true
	}
	return utils.MapKeysToList(groups)
}

// HasEnvironment checks if an environment exists in the project
func (p ProjectConfig) HasEnvironment(name string) bool {
	_, ok := p.Environments[name]
	return ok
}

func (p *ProjectConfig) GetEnvironment(name string) *Environment {
	p.Environments[name].Name = name
	return p.Environments[name]
}

func (p *ProjectConfig) GetStage(env, stage string) *Stage {
	p.GetEnvironment(env).GetStage(stage).Name = stage
	return p.GetEnvironment(env).GetStage(stage)
}

func (p *ProjectConfig) GetCluster(env, stage, cluster string) *Cluster {
	return p.GetStage(env, stage).GetCluster(cluster)
}
