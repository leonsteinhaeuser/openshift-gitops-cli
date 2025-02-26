package menu

import (
	"bufio"
	"fmt"
	"io"
	"path"

	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/utils"
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
		Properties:  properties,
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
		Label: "Select Cluster",
		Items: utils.MapKeysToList(config.Environments[envResult].Stages[stageResult].Clusters),
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "",
			Inactive: "",
			Selected: "",
			Details:  "{{ properties . }}",
			FuncMap: func() map[string]any {
				funcmap := promptui.FuncMap
				funcmap["properties"] = func(clusterName string) string {
					props := config.EnvStageClusterProperty(envResult, stageResult, clusterName)
					resultString := "--------------------------------\nCluster Properties:\n"

					rslt := []string{}
					for k, v := range props {
						rslt = append(rslt, fmt.Sprintf("%s: %s", k, v))
					}

					for _, v := range rslt {
						resultString += "\t" + v + "\n"
					}
					return resultString
				}
				return funcmap
			}(),
		},
	}
	_, clusterResult, err := prompt.Run()
	if err != nil {
		return nil, err
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
		Properties:  properties,
	}, nil
}
