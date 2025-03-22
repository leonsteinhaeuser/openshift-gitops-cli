package main

import (
	"fmt"
	"os"

	"github.com/leonsteinhaeuser/openshift-gitops-cli/internal/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
