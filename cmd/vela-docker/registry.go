// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const (
	credentials = `%s:%s`

	registryFile = `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`

	loginAction = "login"
)

// Registry represents the input parameters for the plugin.
type Registry struct {
	// enable building the image without publishing
	DryRun bool
	// full url to Docker Registry
	Name string
	// password for communication with the Docker Registry
	Password string
	// user name for communication with the Docker Registry
	Username string
}

var (
	// appFs represents a instance of the filesystem.
	appFS = afero.NewOsFs()

	// registryFlags represents for registry settings on the cli.
	registryFlags = []cli.Flag{
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DRY_RUN", "DOCKER_DRY_RUN"},
			FilePath: "/vela/parameters/docker/dry_run,/vela/secrets/docker/dry_run",
			Name:     "registry.dry-run",
			Usage:    "enables building the image without publishing",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REGISTRY", "DOCKER_REGISTRY"},
			FilePath: "/vela/parameters/docker/registry,/vela/secrets/docker/registry",
			Name:     "registry.name",
			Usage:    "Docker registry address to communicate with",
			Value:    "index.docker.io",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: "/vela/parameters/docker/password,/vela/secrets/docker/password",
			Name:     "registry.password",
			Usage:    "password for communication with the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "DOCKER_USERNAME"},
			FilePath: "/vela/parameters/docker/username,/vela/secrets/docker/username",
			Name:     "registry.username",
			Usage:    "user name for communication with the registry",
		},
	}

	// configPath represents the location of the Docker config file for setting registries.
	configPath = "/root/config.json"
)

// Write creates a Docker config.json file for building and publishing the image.
func (r *Registry) Write() error {
	logrus.Trace("creating registry configuration file")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// create basic authentication string for config.json file
	basicAuth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(credentials, r.Username, r.Password)),
	)

	// create output string for config.json file
	out := fmt.Sprintf(
		registryFile,
		r.Name,
		basicAuth,
	)

	return a.WriteFile(configPath, []byte(out), 0644)
}

// Login attempts to authenticate with the registry.
func (r *Registry) Login() error {
	// check if dry run is enabled
	if r.DryRun {
		logrus.Warning("dry_run enabled - skipping authentication with registry")

		// skip authenticating with the registry since dry run is enabled
		return nil
	}

	logrus.Trace("authenticating with registry")

	// variable to store flags for command
	var flags []string

	// add flag for registry password
	flags = append(flags, "--password", r.Password)

	// add flag for registry username
	flags = append(flags, "--username", r.Username)

	// add flag for registry name
	flags = append(flags, r.Name)

	// nolint: gosec // ignore executing command as subprocess
	e := exec.Command(_docker, append([]string{loginAction}, flags...)...)

	// set command stdout to OS stdout
	e.Stdout = os.Stdout
	// set command stderr to OS stderr
	e.Stderr = os.Stderr

	// replace the plain-test password with masking for security purposes
	cmd := strings.ReplaceAll(strings.Join(e.Args, " "), r.Password, constants.SecretMask)

	fmt.Println("$", cmd)

	return e.Run()
}

// Validate verifies the Registry is properly configured.
func (r *Registry) Validate() error {
	logrus.Trace("validating registry plugin configuration")

	// verify registry is provided
	if len(r.Name) == 0 {
		return fmt.Errorf("no registry name provided")
	}

	// check if dry run is disabled
	if !r.DryRun {
		// check if username is provided
		if len(r.Username) == 0 {
			return fmt.Errorf("no registry username provided")
		}

		// check if password is provided
		if len(r.Password) == 0 {
			return fmt.Errorf("no registry password provided")
		}
	}

	return nil
}
