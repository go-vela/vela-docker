// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestDocker_execCmd(t *testing.T) {
	// setup types
	e := exec.Command("echo", "hello")

	err := execCmd(e)
	if err != nil {
		t.Errorf("execCmd returned err: %v", err)
	}
}

func TestDocker_versionCmd(t *testing.T) {
	// setup types
	want := exec.Command(
		_docker,
		"version",
	)

	got := versionCmd()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("versionCmd is %v, want %v", got, want)
	}
}
