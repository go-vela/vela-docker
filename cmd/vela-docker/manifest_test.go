// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"github.com/go-test/deep"
	"os/exec"
	"strings"
	"testing"
)

func TestDocker_Manifest_Create_Command(t *testing.T) {
	// setup types
	// nolint
	m := &Manifest{
		Name: "host.company.com/org/fat-manifest:1",
		Images: []ManifestImage{
			ManifestImage{Name: "host.company.com/org/fat-manifest:1-a"},
			ManifestImage{Name: "host.company.com/org/fat-manifest:1-b"},
		},
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	want := exec.Command(
		_docker,
		manifestAction,
		fmt.Sprintf("create %s", m.Name),
		fmt.Sprintf("--amend %s", m.Images[0].Name),
		fmt.Sprintf("--amend %s", m.Images[1].Name),
	)

	got, _ := m.CreateCommand()
	if !strings.EqualFold(got.String(), want.String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Manifest_WithoutAnnotations(t *testing.T) {
	// setup types
	// nolint
	m := &Manifest{
		Name: "host.company.com/org/fat-manifest:1",
		Images: []ManifestImage{
			ManifestImage{Name: "host.company.com/org/fat-manifest:1-a"},
			ManifestImage{Name: "host.company.com/org/fat-manifest:1-b"},
		},
	}

	got, _ := m.AnnotateCommands()
	if len(got) > 0 {
		t.Errorf("No annotations should have been created")
	}
}

func TestDocker_Manifest_WithVariantAnnotation(t *testing.T) {
	// setup types
	// nolint
	m := &Manifest{
		Name: "host.company.com/org/fat-manifest:1",
		Images: []ManifestImage{
			ManifestImage{
				Name:    "host.company.com/org/fat-manifest:1-a",
				Variant: "v8",
			},
			ManifestImage{Name: "host.company.com/org/fat-manifest:1-b"},
		},
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	want := []*exec.Cmd{
		exec.Command(
			_docker,
			manifestAction,
			"annotate",
			"--variant", m.Images[0].Variant,
			m.Name, m.Images[0].Name,
		),
	}

	got, _ := m.AnnotateCommands()
	if diff := deep.Equal(want, got); diff != nil {
		t.Error(diff)
	}
}
