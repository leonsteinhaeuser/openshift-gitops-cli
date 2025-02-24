package project

type ProjectConfig struct {
	BasePath     string
	Environments map[string]Environment `yaml:"environments"`
}

type Environment struct {
	Name   string           `yaml:"-"`
	Stages map[string]Stage `yaml:"stages"`
}

type Stage struct {
	Name     string             `yaml:"-"`
	Clusters map[string]Cluster `yaml:"clusters"`
}

type Cluster struct {
	Name string `yaml:"-"`
}
