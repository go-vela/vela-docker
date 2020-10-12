// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestDocker_Plugin_Exec(t *testing.T) {
	// TODO Write test
}

func TestDocker_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Context: ".",
			Tags:    []string{"latest"},
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Plugin_Validate_BadBuild(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Context: ".",
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
