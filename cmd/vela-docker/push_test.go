// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestDocker_Push_Command(t *testing.T) {
	// setup types
	// nolint
	p := &Push{
		DisableContentTrust: true,
	}

	want := exec.Command(
		_docker,
		pushAction,
		"--disable-content-trust",
	)

	got, _ := p.Command()
	if !strings.EqualFold(got.String(), want.String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Push_Exec_Error(t *testing.T) {
	// setup types
	p := &Push{}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}
