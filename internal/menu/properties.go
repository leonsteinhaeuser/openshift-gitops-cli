package menu

import (
	"bufio"
	"fmt"
	"io"

	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/cli"
)

func askForProperties(writer io.Writer, reader *bufio.Reader) (map[string]string, error) {
	properties := make(map[string]string)
	for {
		propertyKey, err := cli.StringQuestion(writer, reader, "Property Key", "", func(s string) error {
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
