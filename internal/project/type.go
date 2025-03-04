package project

import (
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/template"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
)

type Environment struct {
	Name       string            `json:"-"`
	Properties map[string]string `json:"properties"`
	Actions    Actions           `json:"actions"`
	Stages     map[string]*Stage `json:"stages"`
}

type Stage struct {
	Name       string              `json:"-"`
	Properties map[string]string   `json:"properties"`
	Actions    Actions             `json:"actions"`
	Clusters   map[string]*Cluster `json:"clusters"`
}

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
	_, ok := p.Environments[env].Stages[stage].Clusters[cluster]
	return ok
}

func (p ProjectConfig) Cluster(env, stage, cluster string) *Cluster {
	return p.Environments[env].Stages[stage].Clusters[cluster]
}

// SetCluster sets the cluster for the given environment and stage
func (p *ProjectConfig) SetCluster(env, stage string, cluster *Cluster) {
	if p.Environments[env].Stages[stage].Clusters == nil {
		p.Environments[env].Stages[stage].Clusters = map[string]*Cluster{}
	}
	p.Environments[env].Stages[stage].Clusters[cluster.Name] = cluster
}

func (p *ProjectConfig) DeleteCluster(env, stage, cluster string) {
	delete(p.Environments[env].Stages[stage].Clusters, cluster)
}

// EnvStageProperty merges the properties of the environment and stage and returns them as a map
func (pc *ProjectConfig) EnvStageProperty(environment, stage string) map[string]string {
	return utils.MergeMaps(pc.Environments[environment].Properties, pc.Environments[environment].Stages[stage].Properties)
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
