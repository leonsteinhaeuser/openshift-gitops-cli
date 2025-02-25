package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/leonsteinhaeuser/openshift-project-bootstrap-cli/internal/project"
)

var (
	projectConfig = &project.ProjectConfig{}
)

const (
	PROJECTFILENAME = "PROJECT.yaml"
)

// check for project file and load it
func init() {
	_, err := os.Stat(PROJECTFILENAME)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println("An error occurred while checking for the PROJECT.yaml file", err)
		os.Exit(1)
		return
	}
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(PROJECTFILENAME)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		f.Close()
	}

	pc, err := project.ParseConfig(PROJECTFILENAME)
	if err != nil {
		fmt.Println("An error occurred while parsing the PROJECT.yaml file", err)
		os.Exit(1)
		return
	}
	projectConfig = pc
}
