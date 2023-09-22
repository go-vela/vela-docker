// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os/exec"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type (
	// Daemon represents the plugin configuration for daemon information.
	Daemon struct {
		// enables specifying a network bridge IP
		Bip string
		// used for translating the storage configuration
		DNS *DNS
		// enables setting custom storage options
		DNSRaw string
		// enable experimental features
		Experimental bool
		// enables insecure registry communication
		InsecureRegistries []string `json:"insecure_registries"`
		// enables IPv6 networking
		IPV6 bool
		// enable setting the log level for the daemon
		LogLevel string `json:"log_level"`
		// enable setting the containers network MTU
		MTU int
		// enables setting a preferred Docker registry mirror
		RegistryMirrors []string `json:"registry_mirrors"`
		// used for translating the storage configuration
		Storage *Storage
		// enables setting custom storage options
		StorageRaw string
	}

	// DNS represents the "dns" prefixed flags within the "dockerd" command.
	DNS struct {
		// enables setting the DNS server to use
		Servers []string
		// enables setting the DNS search domains to use
		Searches []string
	}

	// Storage represents the "storage" prefixed flags within the "dockerd" command.
	Storage struct {
		// enables setting an alternate storage driver
		Driver string
		// enables setting options on the alternate storage driver
		Opts []string
	}
)

// daemonFlags represents for daemon settings on the cli.
var daemonFlags = []cli.Flag{
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_DAEMON", "DOCKER_DAEMON"},
		FilePath: "/vela/parameters/docker/daemon,/vela/secrets/docker/daemon",
		Name:     "daemon",
		Usage:    "enables specifying a network bridge IP",
	},
}

// Command formats and outputs the Build command from
// the provided configuration to push a Docker image.
func (d *Daemon) Command() *exec.Cmd {
	logrus.Trace("creating dockerd command from plugin configuration")

	// variable to store flags for command
	flags := []string{"--data-root=/var/lib/docker", "--host=unix:///var/run/docker.sock"}

	// check if Bip is provided
	if len(d.Bip) > 0 {
		// add flag for Bip from provided build command
		flags = append(flags, "--bip", d.Bip)
	}

	// add flags for DNS configuration
	flags = append(flags, d.DNS.Flags()...)

	// check if Experimental is provided
	if d.Experimental {
		// add flag for Experimental from provided build command
		flags = append(flags, "--experimental")
	}

	// iterate through the insecure registries provided
	for _, i := range d.InsecureRegistries {
		// add flag for InsecureRegistries from provided build command
		flags = append(flags, "--insecure-registry", i)
	}

	// check if Experimental is provided
	if d.IPV6 {
		// add flag for Experimental from provided build command
		flags = append(flags, "--ipv6")
	}

	// check if LogLevel is provided
	if len(d.LogLevel) > 0 {
		// add flag for LogLevel from provided build command
		flags = append(flags, "--log-level", d.LogLevel)
	} else {
		// add flag for LogLevel hardcoded to error level logging
		//
		// this helps to drastically reduce the level of logs
		// output by the plugin when starting up the docker daemon
		flags = append(flags, "--log-level", "error")
	}

	// check if MTU is provided
	if d.MTU > 0 {
		// add flag for MTU from provided build command
		flags = append(flags, "--mtu", strconv.Itoa(d.MTU))
	}

	// iterate through the registry mirrors provided
	for _, r := range d.RegistryMirrors {
		// add flag for RegistryMirrors from provided build command
		flags = append(flags, "--registry-mirror", r)
	}

	// add flags for Storage configuration
	flags = append(flags, d.Storage.Flags()...)

	// the plugin accepts configuration
	return exec.Command(_dockerd, flags...)
}

// Exec formats and runs the commands for pushing a Docker image.
func (d *Daemon) Exec() error {
	logrus.Trace("running dockerd with provided configuration")

	// create the push command for the file
	cmd := d.Command()

	// start the daemon in a thread
	go func() {
		err := execCmd(cmd)
		if err != nil {
			logrus.Error(err)
		}
	}()

	// poll the docker daemon to ensures the daemon is
	// ready to accept connections.
	retryLimit := 5

	// iterate through with a retryLimit
	for i := 0; i < retryLimit; i++ {
		err := versionCmd().Run()
		if err == nil {
			break
		}

		// sleep in between retries
		time.Sleep(time.Duration(i) * time.Second)
	}

	return nil
}

// Flags formats and outputs the flags for
// configuring a Docker daemon DNS settings.
func (d *DNS) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if any dns flags are set
	if d != nil {
		// check if Servers is provided
		if d.Servers != nil {
			for _, d := range d.Servers {
				// add flag for DNS from provided build command
				flags = append(flags, "--dns", d)
			}
		}

		// check if Searches is provided
		if len(d.Searches) > 0 {
			for _, s := range d.Searches {
				// add flag for DNS from provided build command
				flags = append(flags, "--dns-search", s)
			}
		}
	}

	return flags
}

// Flags formats and outputs the flags for
// configuring a Docker daemon DNS settings.
func (s *Storage) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if any storage flags are set
	if s != nil {
		// check if Driver is provided
		if len(s.Driver) > 0 {
			// add flag for Driver from provided build command
			flags = append(flags, "--storage-driver", s.Driver)
		}

		// check if DNSSearch is provided
		if len(s.Opts) > 0 {
			for _, o := range s.Opts {
				// add flag for DNS from provided build command
				flags = append(flags, "--storage-opt", o)
			}
		}
	}

	return flags
}
