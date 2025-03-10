package project

type Stage struct {
	Name       string              `json:"-"`
	Properties map[string]string   `json:"properties"`
	Actions    Actions             `json:"actions"`
	Clusters   map[string]*Cluster `json:"clusters"`
}
