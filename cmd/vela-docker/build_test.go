// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"
)

func TestDocker_Build_Command(t *testing.T) {
	// setup types
	// nolint
	b := &Build{
		AddHosts:     []string{"host.company.com"},
		BuildArgs:    []string{"FOO=BAR"},
		CacheFrom:    "index.docker.in/target/vela-docker",
		CGroupParent: "parent",
		Compress:     true,
		Context:      ".",
		CPU: &CPU{
			Period:  1,
			Quota:   1,
			Shares:  1,
			SetCpus: "(0-3, 0,1)",
			SetMems: "(0-3, 0,1)",
		},
		DisableContentTrust: true,
		File:                "Dockerfile.other",
		ForceRM:             true,
		ImageIDFile:         "path/to/file",
		Isolation:           "hyperv",
		Labels:              []string{"build.number=1", "build.author=octocat"},
		Memory:              []string{"1"},
		MemorySwaps:         []string{"1"},
		Network:             "default",
		NoCache:             true,
		Outputs:             []string{"type=local,dest=path"},
		Platform:            "linux",
		Progress:            "plain",
		Pull:                true,
		Quiet:               true,
		Remove:              true,
		Secrets:             []string{"id=mysecret,src=/local/secret"},
		SecurityOpts:        []string{"seccomp"},
		ShmSizes:            []string{"1"},
		Squash:              true,
		SshComponents:       []string{"default|<id>[=<socket>|<key>[,<key>]]"},
		Stream:              true,
		Tags:                []string{"index.docker.io/target/vela-docker:latest"},
		Target:              "build",
		Ulimits:             []string{"1"},
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	want := exec.Command(
		_docker,
		buildAction,
		fmt.Sprintf("--add-host=%s", b.AddHosts[0]),
		fmt.Sprintf("--build-arg=%s", b.BuildArgs[0]),
		fmt.Sprintf("--cache-from=%s", b.CacheFrom),
		fmt.Sprintf("--cgroup-parent=%s", b.CGroupParent),
		"--compress",
		fmt.Sprintf("--cpu-period=%d", b.CPU.Period),
		fmt.Sprintf("--cpu-quota=%d", b.CPU.Quota),
		fmt.Sprintf("--cpu-shares=%d", b.CPU.Shares),
		fmt.Sprintf("--cpuset-cpus=%s", b.CPU.SetCpus),
		fmt.Sprintf("--cpuset-mems=%s", b.CPU.SetMems),
		"--disable-content-trust",
		fmt.Sprintf("--file=%s", b.File),
		"--force-rm",
		fmt.Sprintf("--iidfile=%s", b.ImageIDFile),
		fmt.Sprintf("--isolation=%s", b.Isolation),
		fmt.Sprintf("--label=%s", b.Labels[0]),
		fmt.Sprintf("--memory=%s", b.Memory[0]),
		fmt.Sprintf("--memory-swap=%s", b.MemorySwaps[0]),
		fmt.Sprintf("--network=%s", b.Network),
		"--no-cache",
		fmt.Sprintf("--output=%s", b.Outputs[0]),
		fmt.Sprintf("--platform=%s", b.Platform),
		fmt.Sprintf("--progress=%s", b.Progress),
		"--pull",
		"--quiet",
		"--rm",
		fmt.Sprintf("--secret=%s", b.Secrets[0]),
		fmt.Sprintf("--security-opt=%s", b.SecurityOpts[0]),
		fmt.Sprintf("--shm-size=%s", b.ShmSizes[0]),
		"--squash",
		fmt.Sprintf("--ssh=%s", b.SshComponents[0]),
		"--stream",
		fmt.Sprintf("--tag=%s", b.Tags[0]),
		fmt.Sprintf("--target=%s", b.Target),
		fmt.Sprintf("--ulimit=%s", b.Ulimits[0]),
		".",
	)

	got, _ := b.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Build_Exec_Error(t *testing.T) {
	// setup types
	b := &Build{
		CPU: &CPU{},
	}

	err := b.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestDocker_Build_Unmarshal(t *testing.T) {
	// setup types
	b := &Build{
		CPURaw: `
  {"period": 1, "quota": 1, "shares": 1, "set_cpus": "(0-3, 0,1)", "set_mems": "(0-3, 0,1)"}
`,
	}

	want := &Build{
		CPU: &CPU{
			Period:  1,
			Quota:   1,
			Shares:  1,
			SetCpus: "(0-3, 0,1)",
			SetMems: "(0-3, 0,1)",
		},
	}

	err := b.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %v", err)
	}

	if !reflect.DeepEqual(b.CPU, want.CPU) {
		t.Errorf("Unmarshal is %v, want %v", b.CPU, want.CPU)
	}
}

func TestDocker_Build_Unmarshal_FailCPUUnmarshal(t *testing.T) {
	// setup types
	b := &Build{
		CPURaw: "!@#$%^&*()",
	}

	err := b.Unmarshal()
	if err == nil {
		t.Errorf("Unmarshal should have returned err")
	}
}

func TestDocker_Build_Validate(t *testing.T) {
	// setup types
	b := &Build{
		Context: ".",
		Tags:    []string{"latest"},
	}

	err := b.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Build_Validate_NoTag(t *testing.T) {
	// setup types
	b := &Build{
		Context: ".",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
