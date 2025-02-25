package menu

import (
	"bufio"
	"fmt"
	"io"
	"path"

	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/project"
	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/utils"
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
	}, nil
}
