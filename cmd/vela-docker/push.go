// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const pushAction = "push"

// Push represents the plugin configuration for push information.
type Push struct {
	// enables naming and optionally a tag in the 'name:tag' format
	Tag string
	// enables skipping image verification (default true)
	DisableContentTrust bool
}

// pushFlags represents for push settings on the cli.
var pushFlags = []cli.Flag{
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_DISABLE_CONTENT_TRUST", "PARAMETER_DISABLE_CONTENT_TRUST"},
		FilePath: "/vela/parameters/docker/disable-content-trust,/vela/secrets/docker/disable-content-trust",
		Name:     "push.disable-content-trust",
		Usage:    "enables skipping image verification (default true)",
	},
}

// Command formats and outputs the Push command from
// the provided configuration to push a Docker image.
func (p *Push) Command() *exec.Cmd {
	logrus.Trace("creating docker push command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if DisableContentTrust is provided
	if p.DisableContentTrust {
		// add flag for DisableContentTrust from provided build command
		flags = append(flags, "--disable-content-trust")
	}

	// add tag to command
	flags = append(flags, p.Tag)

	//nolint: gosec // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_docker, append([]string{pushAction}, flags...)...)
}

// Exec formats and runs the commands for pushing a Docker image.
func (p *Push) Exec() error {
	logrus.Trace("running push with provided configuration")

	// create the push command for the file
	cmd := p.Command()

	// run the push command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}
