package menu

import (
	"bufio"
	"fmt"
	"io"

	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/cli"
	"github.com/leonsteinhaeuser/openshift-gitops-cluster-bootstrap-cli/internal/utils"
	"github.com/manifoldco/promptui"
)

func askForProperties(properties map[string]string, writer io.Writer, reader *bufio.Reader) (map[string]string, error) {
	for {
		prompt := promptui.SelectWithAdd{
			Label:    "Property",
			Items:    utils.MapKeysToList(properties),
			AddLabel: "Create Property",
		}

		_, propertyKey, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if propertyKey == "" {
			return nil, fmt.Errorf("property key cannot be empty")
		}

		// if the key already exists, we ask the user to update the value
		if val, ok := properties[propertyKey]; ok {
			newVal, err := cli.StringQuestion(writer, reader, "Property Value", val, func(s string) error {
				if s == "" {
					return fmt.Errorf("property key cannot be empty")
				}
				if _, ok := properties[s]; ok {
					return fmt.Errorf("property key already exists")
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			properties[propertyKey] = newVal

			isDone, err := cli.BooleanQuestion(writer, reader, "Do you want to add or update another property?", false)
			if err != nil {
				return nil, err
			}
			if !isDone {
				break
			}
			continue
		}

		propertyValue, err := cli.StringQuestion(writer, reader, "Property Value", "", func(s string) error {
			if s == "" {
				return fmt.Errorf("property value cannot be empty")
			}
			if _, ok := properties[s]; ok {
				return fmt.Errorf("property value already exists")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		properties[propertyKey] = propertyValue

		isDone, err := cli.BooleanQuestion(writer, reader, "Do you want to add another property?", false)
		if err != nil {
			return nil, err
		}
		if !isDone {
			break
		}
	}
	return properties, nil
}
