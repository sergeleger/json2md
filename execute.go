// Copyright 2020 Serge Léger. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// +build linux

package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// executeUnsafe executes the command using the default bash interpreter.
func executeUnsafe(cmdStr string) (string, error) {
	errOut := strings.Replace(cmdStr, "\n", "; ", -1)
	log.Println("execute: ", errOut)

	var buf bytes.Buffer
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error: could not execute %q: %w", cmdStr, err)
	}

	return buf.String(), nil
}

// executeSafe never executes the command, it is added to the output wrapped in a code
// block. It is the default mode.
func executeSafe(cmdStr string) (string, error) {
	log.Println("execute: ", cmdStr)

	return "\n```\n" + cmdStr + "\n```\n", nil
}
