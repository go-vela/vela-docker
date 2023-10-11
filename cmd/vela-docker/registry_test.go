// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestDocker_Registry_Login(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		registry *Registry
	}{
		{
			failure: false,
			registry: &Registry{
				Name:     "index.docker.io",
				Username: "octocat",
				Password: "superSecretPassword",
				DryRun:   true,
			},
		},
		{
			failure: true,
			registry: &Registry{
				Name:     "index.docker.io",
				Username: "octocat",
				Password: "superSecretPassword",
				DryRun:   false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.registry.Login()

		if test.failure {
			if err == nil {
				t.Errorf("Login should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Login returned err: %v", err)
		}
	}
}

func TestDocker_Registry_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		registry *Registry
	}{
		{
			failure: false,
			registry: &Registry{
				Name:     "index.docker.io",
				Username: "octocat",
				Password: "superSecretPassword",
				DryRun:   false,
			},
		},
		{
			failure: false,
			registry: &Registry{
				Name:   "index.docker.io",
				DryRun: true,
			},
		},
		{
			failure: true,
			registry: &Registry{
				Username: "octocat",
				Password: "superSecretPassword",
				DryRun:   false,
			},
		},
		{
			failure: true,
			registry: &Registry{
				Name:     "index.docker.io",
				Password: "superSecretPassword",
				DryRun:   false,
			},
		},
		{
			failure: true,
			registry: &Registry{
				Name:     "index.docker.io",
				Username: "octocat",
				DryRun:   false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.registry.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestDocker_Registry_Write(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		Password: "superSecretPassword",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestDocker_Registry_Write_NoName(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Username: "octocat",
		Password: "superSecretPassword",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestDocker_Registry_Write_NoUsername(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestDocker_Registry_Write_NoPassword(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}
