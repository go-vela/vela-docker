// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestDocker_Push_Command(t *testing.T) {
	// setup types
	p := &Push{
		DisableContentTrust: true,
	}

	want := exec.CommandContext(
		t.Context(),
		_docker,
		pushAction,
		"--disable-content-trust ",
	)

	got := p.Command(t.Context())
	if !strings.EqualFold(got.String(), want.String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Push_Exec_Error(t *testing.T) {
	// setup types
	p := &Push{}

	err := p.Exec(t.Context())
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}
