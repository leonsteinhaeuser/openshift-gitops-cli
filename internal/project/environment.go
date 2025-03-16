package project

var (
	_ AddonHandler = &Environment{}
)

type Environment struct {
	Name       string                   `json:"-"`
	Properties map[string]string        `json:"properties"`
	Actions    Actions                  `json:"actions"`
	Stages     map[string]*Stage        `json:"stages"`
	Addons     map[string]*ClusterAddon `json:"addons"`
}

// IsAddonEnabled checks if the addon is enabled for the stage
func (e Environment) IsAddonEnabled(addon string) bool {
	_, ok := e.Addons[addon]
	if !ok {
		return false
	}
	return e.Addons[addon].Enabled
}

// EnableAddon enables the addon for the stage by setting the enabled flag to true
func (e *Environment) EnableAddon(addon string) {
	if e.Addons[addon] == nil {
		e.Addons[addon] = &ClusterAddon{
			Enabled: true,
		}
		return
	}
	e.Addons[addon].Enabled = true
}

// DisableAddon disables the addon for the stage by setting the enabled flag to false
func (e *Environment) DisableAddon(addon string) {
	if _, ok := e.Addons[addon]; !ok {
		// already disabled or not found
		return
	}
	e.Addons[addon].Enabled = false
}

// GetAddons returns the environment addons
func (e *Environment) GetAddons() ClusterAddons {
	return e.Addons
}

// GetAddon returns the addon by name
func (e *Environment) GetAddon(name string) *ClusterAddon {
	return e.Addons[name]
}
