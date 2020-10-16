// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type (
	// Daemon represents the plugin configuration for daemon information.
	// nolint // ignoring length on comments
	Daemon struct {
		// enables specifying a network bridge IP
		Bip string
		// enables a root directory of persistent Docker state (default "/var/lib/docker")
		DataRoot string
		// used for translating the storage configuration
		DNS *DNS
		// enables setting custom storage options
		DNSRaw string
		// enable experimental features
		Experimental bool
		// enables insecure registry communication
		InsecureRegistries []string
		// enables IPv6 networking
		IPV6 bool
		// enable setting the containers network MTU
		MTU int
		// enables setting a preferred Docker registry mirror
		RegistryMirrors []string
		// used for translating the storage configuration
		Storage *Storage
		// enables setting custom storage options
		StorageRaw string
	}

	// DNS represnets the "dns" prefixed flags within the "dockerd" command.
	DNS struct {
		// enables setting the DNS server to use
		Servers []string
		// enables setting the DNS search domains to use
		Searches []string
	}

	// Storage represnets the "storage" prefixed flags within the "dockerd" command.
	Storage struct {
		Driver string
		Opts   []string
	}
)

// daemonFlags represents for daemon settings on the cli.
// nolint // ignoring line length on file paths on comments
var daemonFlags = []cli.Flag{
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_DAEMON"},
		FilePath: string("/vela/parameters/docker/build/daemon,/vela/secrets/docker/build/daemon"),
		Name:     "daemon",
		Usage:    "enables specifying a network bridge IP",
		Value:    true,
	},
}

// Command formats and outputs the Build command from
// the provided configuration to push a Docker image.
// nolint
func (d *Daemon) Command() (*exec.Cmd, error) {
	logrus.Trace("creating dockerd command from plugin configuration")

	// variable to store flags for command
	flags := []string{"--host=unix:///var/run/docker.sock"}

	// check if Bip is provided
	if len(d.Bip) > 0 {
		// add flag for Bip from provided build command
		flags = append(flags, fmt.Sprintf("--bip=%s", d.Bip))
	}

	// check if DataRoot is provided
	if len(d.DataRoot) > 0 {
		// add flag for DataRoot from provided build command
		flags = append(flags, fmt.Sprintf("--data-root=%s", d.DataRoot))
	}

	// add flags for DNS configuration
	flags = append(flags, d.DNS.Flags()...)

	// check if Experimental is provided
	if d.Experimental {
		// add flag for Experimental from provided build command
		flags = append(flags, "--experimental")
	}

	// check if InsecureRegistries is provided
	if len(d.InsecureRegistries) > 0 {
		var args string
		for _, arg := range d.InsecureRegistries {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for InsecureRegistries from provided build command
		flags = append(flags, fmt.Sprintf("--insecure-registry=%s", strings.TrimPrefix(args, " ")))
	}

	// check if Experimental is provided
	if d.IPV6 {
		// add flag for Experimental from provided build command
		flags = append(flags, "--ipv6")
	}

	// check if MTU is provided
	if d.MTU > 0 {
		// add flag for MTU from provided build command
		flags = append(flags, fmt.Sprintf("--mtu=%d", d.MTU))
	}

	// check if RegistryMirrors is provided
	if len(d.RegistryMirrors) > 0 {
		var args string
		for _, arg := range d.RegistryMirrors {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for RegistryMirrors from provided build command
		flags = append(flags, fmt.Sprintf("--registry-mirror=%s", strings.TrimPrefix(args, " ")))
	}

	// add flags for Storage configuration
	flags = append(flags, d.Storage.Flags()...)

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_dockerd, flags...), nil
}

// Exec formats and runs the commands for pushing a Docker image.
func (d *Daemon) Exec() error {
	logrus.Trace("running dockerd with provided configuration")

	// create the push command for the file
	cmd, err := d.Command()
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

// Flags formats and outputs the flags for
// configuring a Docker daemon DNS settings.
func (d *DNS) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if Servers is provided
	if len(d.Servers) > 0 {
		var args string
		for _, arg := range d.Servers {
			args += fmt.Sprintf("--dns=%s ", strings.TrimPrefix(arg, " "))
		}
		// add flag for DNS from provided build command
		flags = append(flags, args)
	}

	// check if Searches is provided
	if len(d.Searches) > 0 {
		var args string
		for _, arg := range d.Searches {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for DNS from provided build command
		flags = append(flags, fmt.Sprintf("--dns-search=%s", strings.TrimPrefix(args, " ")))
	}

	return flags
}

// Flags formats and outputs the flags for
// configuring a Docker daemon DNS settings.
func (s *Storage) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if Driver is provided
	if len(s.Driver) > 0 {
		// add flag for Driver from provided build command
		flags = append(flags, fmt.Sprintf("--storage-driver=%s", s.Driver))
	}

	// check if DNSSearch is provided
	if len(s.Opts) > 0 {
		var args string
		for _, arg := range s.Opts {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for DNS from provided build command
		flags = append(flags, fmt.Sprintf("--storage-opt=%s", strings.TrimPrefix(args, " ")))
	}

	return flags
}
