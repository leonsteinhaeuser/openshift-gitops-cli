package project

type Environment struct {
	Name       string            `json:"-"`
	Properties map[string]string `json:"properties"`
	Actions    Actions           `json:"actions"`
	Stages     map[string]*Stage `json:"stages"`
}
