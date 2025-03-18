package project

var (
	_ AddonHandler = &Stage{}
)

type Stage struct {
	Name       string                   `json:"-"`
	Properties map[string]string        `json:"properties"`
	Actions    Actions                  `json:"actions"`
	Clusters   map[string]*Cluster      `json:"clusters"`
	Addons     map[string]*ClusterAddon `json:"addons"`
}

// IsAddonEnabled checks if the addon is enabled for the stage
func (s Stage) IsAddonEnabled(addon string) bool {
	_, ok := s.Addons[addon]
	if !ok {
		return false
	}
	return s.Addons[addon].Enabled
}

// EnableAddon enables the addon for the stage by setting the enabled flag to true
func (s *Stage) EnableAddon(addon string) {
	if s.Addons[addon] == nil {
		s.Addons[addon] = &ClusterAddon{
			Enabled: true,
		}
		return
	}
	s.Addons[addon].Enabled = true
}

// DisableAddon disables the addon for the stage by setting the enabled flag to false
func (s *Stage) DisableAddon(addon string) {
	if _, ok := s.Addons[addon]; !ok {
		// already disabled or not found
		return
	}
	s.Addons[addon].Enabled = false
}

// GetAddons returns the environment addons
func (s *Stage) GetAddons() ClusterAddons {
	return s.Addons
}

// GetAddon returns the addon by name
func (s *Stage) GetAddon(name string) *ClusterAddon {
	return s.Addons[name]
}

// GetCluster returns the cluster by name
func (s *Stage) GetCluster(name string) *Cluster {
	return s.Clusters[name]
}
