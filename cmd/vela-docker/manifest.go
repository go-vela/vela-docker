// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

// nolint
const manifestAction = "manifest"

type (
	// Manifest represents the plugin configuration for build information.
	// nolint // ignoring length on comments
	Manifest struct {
		// name of the fat manifest to create
		Name string `json:"name"`
		// list of existing images to include in the created fat manifest
		Images     []ManifestImage `json:"images"`
		ShouldExec bool            `json:"-"`
		RawSpec    string          `json:"-"`
	}

	// ManifestImage represents the component image as well as the annotations available
	// in the "docker manifest annotate" command.
	ManifestImage struct {
		// name of the image to be annotated
		Name string `json:"name"`
		// enables setting the --arch annotation
		Arch string `json:"arch"`
		// enables setting the --os annotation
		OS string `json:"os"`
		// enables setting the --os-features annotation
		OSFeatures []string `json:"os_features"`
		// enables setting the --os-version annotation
		OSVersion string `json:"os_version"`
		// enables setting the --variant annotation
		Variant string `json:"variant"`
	}
)

// manifestFlags represents for build settings on the cli.
// nolint // ignoring line length on file paths on comments
var manifestFlags = []cli.Flag{
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_MANIFEST", "DOCKER_MANIFEST"},
		FilePath: "/vela/parameters/docker/manifest,vela/secrets/docker/manifest",
		Name:     "manifest.spec",
		Usage:    "manifest ",
	},
}

func (m *Manifest) Unmarshal() error {
	if m == nil || len(m.RawSpec) == 0 {
		if m != nil {
			m.ShouldExec = false
		}
		return fmt.Errorf("No manifest provided")
	}
	logrus.Tracef("Parsing manifest.spec %s", m.RawSpec)
	err := json.Unmarshal([]byte(m.RawSpec), &m)
	if err != nil {
		m.ShouldExec = false
		return err
	}
	return nil
}

// CreateCommand formats and outputs the Manifest.Create command from
// the provided configuration to build a Docker fat manifest.
// nolint
func (m *Manifest) CreateCommand() (*exec.Cmd, error) {
	logrus.Trace("creating docker manifest create command from plugin configuration")

	// variable to store flags for command
	var flags []string

	flags = append(flags, "create", m.Name)

	// iterate through the amended images provided
	for _, image := range m.Images {
		// add flag for amended images from provided manifest create command
		flags = append(flags, "--amend", image.Name)
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_docker, append([]string{manifestAction}, flags...)...), nil
}

// CreateCommand formats and outputs the Manifest.Create command from
// the provided configuration to build a Docker fat manifest.
// nolint
func (m *Manifest) AnnotateCommands() ([]*exec.Cmd, error) {
	logrus.Trace("creating docker manifest annotate command from plugin configuration")

	// variable to store flags for command
	annotateCommands := []*exec.Cmd{}

	// iterate through the amended images first to see if we have any required annotations
	for _, image := range m.Images {
		var imageFlags []string

		if len(image.Arch) > 0 {
			imageFlags = append(imageFlags, "--arch", image.Arch)
		}
		if len(image.OS) > 0 {
			imageFlags = append(imageFlags, "--os", image.OS)
		}
		if len(image.OSFeatures) > 0 {
			imageFlags = append(imageFlags, "--os-features", strings.Join(image.OSFeatures, ","))
		}
		if len(image.OSVersion) > 0 {
			imageFlags = append(imageFlags, "--os-version", image.OSVersion)
		}
		if len(image.Variant) > 0 {
			imageFlags = append(imageFlags, "--variant", image.Variant)
		}

		if len(imageFlags) > 0 {
			logrus.Tracef("Annotating %s in list %s with %s", image.Name, m.Name, strings.Join(imageFlags, " "))
			cmdFlags := append([]string{manifestAction}, "annotate")
			cmdFlags = append(cmdFlags, imageFlags...)
			cmdFlags = append(cmdFlags, m.Name, image.Name)
			annotateCommands = append(annotateCommands,
				exec.Command(_docker, cmdFlags...))

		}
	}
	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return annotateCommands, nil
}

// CreateCommand formats and outputs the Manifest.Create command from
// the provided configuration to build a Docker fat manifest.
// nolint
func (m *Manifest) PushCommand() (*exec.Cmd, error) {
	logrus.Trace("creating docker manifest push command from plugin configuration")

	// variable to store flags for command
	flags := []string{"push", m.Name}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_docker, append([]string{manifestAction}, flags...)...), nil
}

func (m *Manifest) execCreateCmd() error {
	c, err := m.CreateCommand()
	if err != nil {
		return err
	}
	return execCmd(c)
}

func (m *Manifest) execAnnotateCmds() error {
	cmds, err := m.AnnotateCommands()
	if err != nil {
		return err
	}
	for _, c := range cmds {
		err = execCmd(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manifest) ExecPushCmd() error {
	if nil == m || ! m.ShouldExec {
		return fmt.Errorf("Cannot push an invalid manifest")
	}
	c, err := m.PushCommand()
	if err != nil {
		return err
	}
	return execCmd(c)
}

// Exec formats and runs the commands for building a Docker image.
func (m *Manifest) Exec() error {
	if m == nil || !m.ShouldExec {
		return fmt.Errorf("Manifest step not executing because manifest was invalid")
	}
	logrus.Trace("running manifest with provided configuration")

	if len(m.Name) > 0 {
		err := m.execCreateCmd()
		if err != nil {
			return err
		}
		err = m.execAnnotateCmds()
		if err != nil {
			return err
		}
	}

	// run the build command for the file

	return nil
}

// Validate verifies the Build is properly configured.
func (m *Manifest) Validate() error {
	logrus.Trace("validating manifest plugin configuration")

	if len(m.Name) == 0 {
		return fmt.Errorf("cannot create a fat manifest without a name")
	}

	if len(m.Name) > 0 && len(m.Images) == 0 {
		m.ShouldExec = false
		return fmt.Errorf("cannot create a fat manifest without including the images to be included")
	}
	m.ShouldExec = true
	return nil
}
