// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// _docker is the path to the executable binary in the image.
	_docker = "/usr/local/bin/docker"
	// _dockerd is the path to the executable daemon in the image.
	_dockerd = "/usr/local/bin/rootlesskit"
)

// execCmd is a helper function to
// run the provided command.
func execCmd(e *exec.Cmd) error {
	logrus.Tracef("executing cmd %s", strings.Join(e.Args, " "))

	// set command stdout to OS stdout
	e.Stdout = os.Stdout
	// set command stderr to OS stderr
	e.Stderr = os.Stderr

	// output "trace" string for command
	fmt.Println("$", strings.Join(e.Args, " "))

	return e.Run()
}

// infoCmd is a helper function to check if
// the daemon is ready.
func infoCmd() *exec.Cmd {
	logrus.Trace("creating docker info command")

	// variable to store flags for command
	var flags []string

	// add flag for version img command
	flags = append(flags, "info")

	return exec.Command(_docker, flags...)
}

// versionCmd is a helper function to output
// the client and server version information.
func versionCmd() *exec.Cmd {
	logrus.Trace("creating docker version command")

	// variable to store flags for command
	var flags []string

	// add flag for version img command
	flags = append(flags, "version")

	return exec.Command(_docker, flags...)
}
