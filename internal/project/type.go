package project

type ProjectConfig struct {
	BasePath         string                 `json:"basePath"`
	TemplateBasePath string                 `json:"templateBasePath"`
	Environments     map[string]Environment `json:"environments"`
}

type Environment struct {
	Name   string           `json:"-"`
	Stages map[string]Stage `json:"stages"`
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
