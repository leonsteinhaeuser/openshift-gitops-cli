package utils

import (
	"io"
	"os/exec"
)

// ExecuteShellCommand executes a shell command with the given arguments and writes the output to the given writer
func ExecuteShellCommand(stdout, errout io.Writer, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = errout
	return cmd.Run()
}
