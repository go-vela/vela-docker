// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// build arguments loaded for the plugin
	Build *Build
	// daemon arguments loaded for the plugin
	Daemon *Daemon
	// push arguments loaded for the plugin
	Push *Push
	// registry arguments loaded for the plugin
	Registry *Registry
}

// Exec formats and runs the commands for building and publishing a Docker image.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// start the docker daemon with configuration
	err := p.Daemon.Exec()
	if err != nil {
		return err
	}

	// create registry file for authentication
	err = p.Registry.Write()
	if err != nil {
		return err
	}

	// create registry login to validate authentication
	err = p.Registry.Login()
	if err != nil {
		return err
	}

	// execute build configuration
	err = p.Build.Exec()
	if err != nil {
		return err
	}

	// check if registry dry run is enabled
	if !p.Registry.DryRun {
		// push all tags
		for _, t := range p.Build.Tags {
			// set the tag to be pushed
			p.Push.Tag = t

			// execute push configuration
			err = p.Push.Exec()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate verifies the Plugin is properly configured.
func (p *Plugin) Validate(daemon string) error {
	logrus.Debug("validating plugin configuration")

	// serialize daemon settings into plugin
	if len(daemon) > 0 {
		err := json.Unmarshal([]byte(daemon), &p.Daemon)
		if err != nil {
			return err
		}
	}

	// validate registry configuration
	err := p.Registry.Validate()
	if err != nil {
		return err
	}

	// when user adds configuration additional options
	err = p.Build.Unmarshal()
	if err != nil {
		return err
	}

	// validate build configuration
	err = p.Build.Validate()
	if err != nil {
		return err
	}

	return nil
}
