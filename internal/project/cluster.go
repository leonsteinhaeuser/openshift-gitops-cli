package project

import "fmt"

type ClusterAddon struct {
	Enabled    bool           `json:"enabled"`
	Properties map[string]any `json:"properties"`
}

type Cluster struct {
	Name       string                   `json:"-"`
	Addons     map[string]*ClusterAddon `json:"addons"`
	Properties map[string]string        `json:"properties"`
}

// IsAddonEnabled checks if the addon is enabled for the cluster
func (c Cluster) IsAddonEnabled(addon string) bool {
	_, ok := c.Addons[addon]
	if !ok {
		return false
	}
	return c.Addons[addon].Enabled
}

// EnableAddon enables the addon for the cluster by setting the enabled flag to true
func (c *Cluster) EnableAddon(addon string) {
	if c.Addons[addon] == nil {
		c.Addons[addon] = &ClusterAddon{
			Enabled: true,
		}
		return
	}
	c.Addons[addon].Enabled = true
	fmt.Println("Enabled addon option:", c.Addons[addon].Enabled)
}

// DisableAddon disables the addon for the cluster by setting the enabled flag to false and removing all properties
func (c *Cluster) DisableAddon(addon string) {
	if _, ok := c.Addons[addon]; !ok {
		// already disabled or not found
		return
	}
	c.Addons[addon].Enabled = false
	c.Addons[addon].Properties = map[string]any{}
}
