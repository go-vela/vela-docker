// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestDocker_Daemon_Command(t *testing.T) {
	// setup types
	d := &Daemon{
		Bip: "192.168.1.5/24",
		DNS: &DNS{
			Servers:  []string{"10.20.1.2", "10.20.1.3"},
			Searches: []string{"8.8.8.8"},
		},
		Experimental:       true,
		InsecureRegistries: []string{"private.registry.com"},
		IPV6:               true,
		MTU:                1500,
		RegistryMirrors:    []string{"mirror.registry.com"},
		Storage: &Storage{
			Driver: "overlay2",
			Opts:   []string{"ftype=1"},
		},
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	want := exec.Command(
		_dockerd,
		"--data-root=/var/lib/docker",
		"--host=unix:///var/run/docker.sock",
		fmt.Sprintf("--bip %s", d.Bip),
		fmt.Sprintf("--dns %s --dns %s", d.DNS.Servers[0], d.DNS.Servers[1]),
		fmt.Sprintf("--dns-search %s", d.DNS.Searches[0]),
		"--experimental",
		fmt.Sprintf("--insecure-registry %s", d.InsecureRegistries[0]),
		"--ipv6",
		fmt.Sprintf("--mtu %d", d.MTU),
		fmt.Sprintf("--registry-mirror %s", d.RegistryMirrors[0]),
		fmt.Sprintf("--storage-driver %s", d.Storage.Driver),
		fmt.Sprintf("--storage-opt %s", d.Storage.Opts[0]),
	)

	got, _ := d.Command()
	if !strings.EqualFold(got.String(), want.String()) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Daemon_Exec_Error(t *testing.T) {
	// setup types
	d := &Daemon{
		DNS:     &DNS{},
		Storage: &Storage{},
	}

	err := d.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}
