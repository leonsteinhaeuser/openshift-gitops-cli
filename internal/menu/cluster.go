package menu

import (
	"bufio"
	"fmt"
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

type clusterMenu struct {
	writer io.Writer
	reader *bufio.Reader
	config *project.ProjectConfig
}

// menuCreateCluster creates a context menu to create a new cluster
func (c *clusterMenu) menuCreateCluster(env, stage string) (*project.Cluster, error) {
	clusterName, err := cli.StringQuestion(c.writer, c.reader, "Cluster Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("cluster name cannot be empty")
		}
		if c.config.HasCluster(env, stage, s) {
			return fmt.Errorf("cluster already exists")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	cluster := &project.Cluster{
		Name:       clusterName,
		Addons:     map[string]*project.ClusterAddon{},
		Properties: map[string]string{},
	}
	cluster.SetDefaultAddons(c.config)

	err = c.menuSettings(env, stage, cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// menuSettings creates a context menu to manage the settings of a cluster
func (c *clusterMenu) menuSettings(env, stage string, cluster *project.Cluster) error {
	for {
		prompt := promptui.Select{
			Label: "Settings",
			Items: []string{"Addons", "Properties", "Done"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		switch result {
		case "Addons":
			addon := addonClusterMenu{
				writer: c.writer,
				reader: c.reader,
				config: c.config,
			}

			err := addon.menuManageAddons(cluster)
			if err != nil {
				return err
			}
		case "Properties":
			properties, err := c.menuClusterSettingsProperties(env, stage, cluster)
			if err != nil {
				return err
			}
			cluster.Properties = properties
		case "Done":
			return nil
		default:
			return fmt.Errorf("invalid option %s", result)
		}
	}
}

// menuUpdateCluster creates a context menu to update an existing cluster
func (c *clusterMenu) menuUpdateCluster(envName, stageName, clusterName string) (*project.Cluster, error) {
	cluster := c.config.Cluster(envName, stageName, clusterName)
	if cluster.Name == "" {
		cluster.Name = clusterName
	}
	cluster.SetDefaultAddons(c.config)
	err := c.menuSettings(envName, stageName, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// menuDeleteCluster creates a context menu to delete an existing cluster
func (c *clusterMenu) menuDeleteCluster(env, stage, cluster string) (*project.Cluster, error) {
	confirmation, err := cli.BooleanQuestion(c.writer, c.reader, "Are you sure to delete the cluster?", false)
	if err != nil {
		return nil, err
	}
	if !confirmation {
		return nil, fmt.Errorf("confirmation denied")
	}
	return c.config.Cluster(env, stage, cluster), nil
}

func (c *clusterMenu) menuClusterSettingsProperties(env, stage string, cluster *project.Cluster) (map[string]string, error) {
	clusterProperties := map[string]string{}
	for {
		properties := utils.MergeMaps(c.config.EnvStageProperty(env, stage), cluster.Properties, clusterProperties)

		prompt := promptui.SelectWithAdd{
			Label:    "Properties",
			Items:    append(utils.SortStringSlice(utils.MapKeysToList(properties)), "Done"),
			AddLabel: "Create Property",
		}
		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		if result == "" {
			return nil, fmt.Errorf("property key cannot be empty")
		}
		if result == "Done" {
			// user is done
			break
		}

		val, err := cli.StringQuestion(c.writer, c.reader, "Property Value", properties[result], func(s string) error {
			if s == "" {
				return fmt.Errorf("property value cannot be empty")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		clusterProperties[result] = val
	}
	return clusterProperties, nil
}
