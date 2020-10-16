// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// nolint
const pushAction = "push"

// Push represents the plugin configuration for push information.
type Push struct {
	// enables skipping image verification (default true)
	DisableContentTrust bool
}

// pushFlags represents for push settings on the cli.
// nolint // ignoring line length on file paths on comments
var pushFlags = []cli.Flag{
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_DISABLE_CONTENT_TRUST"},
		FilePath: string("/vela/parameters/docker/build/disable-content-trust,/vela/secrets/docker/build/disable-content-trust"),
		Name:     "push.disable-content-trust",
		Usage:    "enables skipping image verification (default true)",
		Value:    true,
	},
}

// Command formats and outputs the Push command from
// the provided configuration to push a Docker image.
// nolint
func (p *Push) Command() (*exec.Cmd, error) {
	logrus.Trace("creating docker push command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if DisableContentTrust is provided
	if p.DisableContentTrust {
		// add flag for DisableContentTrust from provided build command
		flags = append(flags, "--disable-content-trust")
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_docker, append([]string{pushAction}, flags...)...), nil
}

// Exec formats and runs the commands for pushing a Docker image.
func (p *Push) Exec() error {
	logrus.Trace("running push with provided configuration")

	// create the push command for the file
	cmd, err := p.Command()
	if err != nil {
		return err
	}

	// run the push command for the file
	err = execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}
