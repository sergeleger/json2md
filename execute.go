// +build withExecute

package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func init() {
	templateFuncs["execute"] = execute
}

// execute executes the command using the default bash interpreter.
func execute(cmdStr string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error: count not execute %q: %w", cmdStr, err)
	}

	return buf.String(), nil
}
