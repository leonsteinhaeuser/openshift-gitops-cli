package menu

import (
	"bufio"
	"fmt"
	"io"
	"path"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

// CreateCluster creates a context menu to create a new cluster
// As part of this we ask for the environment, stage and cluster name
func CreateCluster(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) (*CarrierCreateCluster, error) {
	// read environments
	prompt := promptui.Select{
		Label: "Select Environment",
		Items: utils.MapKeysToList(config.Environments),
	}
	_, envResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	// read stages
	stages := config.Environments[envResult].Stages
	prompt = promptui.Select{
		Label: "Select Stage",
		Items: utils.MapKeysToList(stages),
	}
	_, stageResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	// ask for cluster name
	clusterName, err := cli.StringQuestion(writer, reader, "Cluster Name", "", func(s string) error {
		if s == "" {
			return fmt.Errorf("cluster name cannot be empty")
		}

		if _, ok := stages[stageResult].Clusters[s]; ok {
			return fmt.Errorf("cluster already exists")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	addonConfig := map[string]map[string]any{}
	for addonName, cfg := range config.ParsedAddons {
		shouldAddonBeEnabled, err := cli.BooleanQuestion(writer, reader, fmt.Sprintf("Do you want to enable %s?", addonName), config.Addons[addonName].DefaultEnabled)
		if err != nil {
			return nil, err
		}
		if !shouldAddonBeEnabled {
			// addon should not be enabled
			continue
		}

		// ask for addon properties
		addonProperties, err := templateManifestPropertiesMenu(
			writer,
			reader,
			config.Environments[envResult].Stages[stageResult].Clusters[clusterName],
			cfg)
		if err != nil {
			return nil, err
		}
		addonConfig[addonName] = addonProperties
	}

	// let's ask if the user want to add additional properties
	createProperties, err := cli.BooleanQuestion(writer, reader, "Do you want to add properties?", false)
	if err != nil {
		return nil, err
	}
	properties := map[string]string{}
	if createProperties {
		pts, err := askForProperties(config.EnvStageProperty(envResult, stageResult), writer, reader)
		if err != nil {
			return nil, err
		}
		properties = pts
	}

	// ask for confirmation
	fqnPath := path.Join(config.BasePath, envResult, stageResult, clusterName)
	confirmation, err := cli.BooleanQuestion(writer, reader, fmt.Sprintf("Are you sure to create a new cluster in %s", fqnPath), false)
	if err != nil {
		return nil, err
	}
	if !confirmation {
		return nil, fmt.Errorf("confirmation denied")
	}

	return &CarrierCreateCluster{
		Environment: envResult,
		Stage:       stageResult,
		ClusterName: clusterName,
		Addons:      addonConfig,
		Properties:  utils.ReduceMap(properties, config.EnvStageProperty(envResult, stageResult)),
	}, nil
}

// UpdateCluster creates a context menu to update an existing cluster
// As part of this we ask for the environment, stage and cluster
func UpdateCluster(config *project.ProjectConfig, writer io.Writer, reader *bufio.Reader) (*CarrierCreateCluster, error) {
	// read environments
	prompt := promptui.Select{
		Label: "Select Environment",
		Items: utils.MapKeysToList(config.Environments),
	}
	_, envResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	// read stages
	prompt = promptui.Select{
		Label: "Select Stage",
		Items: utils.MapKeysToList(config.Environments[envResult].Stages),
	}
	_, stageResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	// read clusters
	prompt = promptui.Select{
		Label:     "Select Cluster",
		Items:     utils.MapKeysToList(config.Environments[envResult].Stages[stageResult].Clusters),
		Templates: helperSelectTemplate(config, envResult, stageResult),
	}
	_, clusterResult, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	addonConfig := map[string]map[string]any{}
	for addonName, cfg := range config.ParsedAddons {
		shouldAddonBeEnabled, err := cli.BooleanQuestion(writer, reader, fmt.Sprintf("Do you want to enable %s?", addonName), config.Addons[addonName].DefaultEnabled)
		if err != nil {
			return nil, err
		}
		if !shouldAddonBeEnabled {
			// addon should not be enabled
			continue
		}

		// ask for addon properties
		addonProperties, err := templateManifestPropertiesMenu(
			writer,
			reader,
			config.Environments[envResult].Stages[stageResult].Clusters[clusterResult],
			cfg)
		if err != nil {
			return nil, err
		}
		addonConfig[addonName] = addonProperties
	}

	// let's ask if the user want to add additional properties
	createProperties, err := cli.BooleanQuestion(writer, reader, "Do you want to update properties?", false)
	if err != nil {
		return nil, err
	}
	properties := map[string]string{}
	if createProperties {
		pts, err := askForProperties(config.EnvStageClusterProperty(envResult, stageResult, clusterResult), writer, reader)
		if err != nil {
			return nil, err
		}
		properties = pts
	}

	return &CarrierCreateCluster{
		Environment: envResult,
		Stage:       stageResult,
		ClusterName: clusterResult,
		Addons:      addonConfig,
		Properties:  utils.ReduceMap(properties, config.EnvStageProperty(envResult, stageResult)),
	}, nil
}
