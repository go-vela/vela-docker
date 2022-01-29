// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
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
	// manifest arguments loaded for the plugin
	Manifest *Manifest
}

// Exec formats and runs the commands for building and publishing a Docker image.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// start the docker daemon with configuration
	err := p.Daemon.Exec()
	if err != nil {
		return err
	}

	// output the docker version
	err = execCmd(versionCmd())
	if err != nil {
		return err
	}

	// output the docker information
	err = execCmd(infoCmd())
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
	if p.Build.ShouldExec {
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
	}


	if p.Manifest.ShouldExec {
		err = p.Manifest.Exec()
		if err != nil {
			return err
		}
		if !p.Registry.DryRun {
			err = p.Manifest.ExecPushCmd()
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
	buildErr := p.Build.Unmarshal()
	if buildErr != nil {
		logrus.Tracef("Problem unmarshalling build: %v", buildErr)
	} else {
		// validate build configuration
		buildErr = p.Build.Validate()
		if buildErr != nil {
			logrus.Tracef("Problem validating build: %v", buildErr)
		}
	}

	manifestErr := p.Manifest.Unmarshal()
	if manifestErr != nil {
		logrus.Tracef("Problem unmarshalling manifest: %v", manifestErr)
	} else {
		manifestErr = p.Manifest.Validate()
		if manifestErr != nil {
			logrus.Tracef("Problem validating manifest: %v", manifestErr)
		}
	}

	if manifestErr != nil && buildErr != nil {
		err = fmt.Errorf("Neither build nor manifest was valid")
		return err
	}
	return nil
}
