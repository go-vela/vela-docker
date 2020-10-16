// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// nolint
const buildAction = "build"

type (
	// Build represents the plugin configuration for build information.
	// nolint // ignoring length on comments
	Build struct {
		// enables adding a custom host-to-IP mapping (host:ip)
		AddHosts []string
		// enables setting build-time variables
		BuildArgs []string
		// enables setting images to consider as cache sources
		CacheFrom string
		// enables setting an optional parent cgroup for the container
		CGroupParent string
		// enables setting compression the build context using gzip
		Compress bool
		// enables setting the build context
		Context string
		// used for translating the cpu configuration
		CPU *CPU
		// enables setting custom cpu options
		CPURaw string
		// enables skipping image verification (default true)
		DisableContentTrust bool
		// enables setting the name of the Dockerfile (Default is 'PATH/Dockerfile')
		File string
		// enables setting always remove on intermediate containers
		ForceRM bool
		// enables setting writing the image ID to the file
		ImageIDFile string
		// enables container isolation technology
		Isolation string
		// used for translating the pre-defined image labels
		Label *Label
		// enables setting metadata for an image
		Labels []string
		// enables setting a memory limit
		Memory []string
		// enables setting a swap limit equal to memory plus swap: '-1' to enable unlimited swap
		MemorySwaps []string
		// enables setting the networking mode for the RUN instructions during build (default "default")
		Network string
		// enables setting not use cache when building the image
		NoCache bool
		// enables setting an output destination (format: type=local,dest=path)
		Outputs []string
		// enables setting a platform if server is multi-platform capable
		Platform string
		// enables setting type of progress output - options (auto|plain|tty)
		Progress string
		// enables always attempting to pull a newer version of the image
		Pull bool
		// enables suppressing the build output and print image ID on success
		Quiet bool
		// enables removing the intermediate containers after a successful build (default true)
		Remove bool
		// enables setting a secret file to expose to the build (only if BuildKit enabled): id=mysecret,src=/local/secret
		Secrets []string
		// enables setting security options
		SecurityOpts []string
		// enables setting the size of /dev/shm
		ShmSizes []string
		// enables setting squash newly built layers into a single new layer
		Squash bool
		// enables setting an ssh agent socket or keys to expose to the build (only if BuildKit enabled) (format: default|<id>[=<socket>|<key>[,<key>]])
		SshComponents []string
		// enables streaming attaches to server to negotiate build context
		Stream bool
		// enables naming and optionally a tag in the 'name:tag' format
		Tags []string
		// enables setting the target build stage to build.
		Target string
		// enables setting ulimit options (default [])
		Ulimits []string
	}

	// CPU represents the "cpu" prefixed flags within the "docker build" command.
	CPU struct {
		// enables setting limits on the CPU CFS (Completely Fair Scheduler) period
		Period int
		// enables setting limit on the CPU CFS (Completely Fair Scheduler) quota
		Quota int
		// enables setting CPU shares (relative weight)
		Shares int
		// enables setting CPUs in which to allow execution (0-3, 0,1)
		SetCpus string `json:"set_cpus"`
		// enables setting MEMs in which to allow execution (0-3, 0,1)
		SetMems string `json:"set_mems"`
	}

	// Label represents the open image specification fields.
	Label struct {
		// author from the source commit
		AuthorEmail string
		// commit sha from the source commit
		Commit string
		// timestamp when the image was built
		Created string
		// full name of the repository
		FullName string
		// build number from vela
		Number int
		// direct url of the repository
		URL string
	}
)

// buildFlags represents for build settings on the cli.
// nolint // ignoring line length on file paths on comments
var buildFlags = []cli.Flag{
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_ADD_HOSTS"},
		FilePath: string("/vela/parameters/docker/build/add_hosts,/vela/secrets/docker/build/add_hosts"),
		Name:     "build.add-hosts",
		Usage:    "enables adding a custom host-to-IP mapping (host:ip)",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_BUILD_ARGS"},
		FilePath: string("/vela/parameters/docker/build/build_args,/vela/secrets/docker/build/build_args"),
		Name:     "build.build-args",
		Usage:    "enables setting build time arguments for the dockerfile",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_CACHE_FROM"},
		FilePath: string("/vela/parameters/docker/build/cache_from,/vela/secrets/docker/build/cache_from"),
		Name:     "build.cache-from",
		Usage:    "enables setting images to consider as cache sources",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_CGROUP_PARENT"},
		FilePath: string("/vela/parameters/docker/build/cgroup_parent,/vela/secrets/docker/build/cgroup_parent"),
		Name:     "build.cgroup-parent",
		Usage:    "enables setting an optional parent cgroup for the container",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_COMPRESS"},
		FilePath: string("/vela/parameters/docker/build/compress,/vela/secrets/docker/build/compress"),
		Name:     "build.compress",
		Usage:    "enables setting compression the build context using gzip",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_CONTEXT"},
		FilePath: string("/vela/parameters/docker/build/context,/vela/secrets/docker/build/context"),
		Name:     "build.context",
		Usage:    "enables setting the build context",
		Value:    ".",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_CPU"},
		FilePath: string("/vela/parameters/docker/build/cpu,/vela/secrets/docker/build/cpu"),
		Name:     "build.cpu",
		Usage:    "enables setting custom cpu options",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_DISABLE_CONTENT_TRUST"},
		FilePath: string("/vela/parameters/docker/build/disable-content-trust,/vela/secrets/docker/build/disable-content-trust"),
		Name:     "build.disable-content-trust",
		Usage:    "enables skipping image verification (default true)",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_FILE"},
		FilePath: string("/vela/parameters/docker/build/file,/vela/secrets/docker/build/file"),
		Name:     "build.file",
		Usage:    "enables setting the name of the Dockerfile (Default is 'PATH/Dockerfile')",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_FORCE_RM"},
		FilePath: string("/vela/parameters/docker/build/force_rm,/vela/secrets/docker/build/force_rm"),
		Name:     "build.force-rm",
		Usage:    "enables setting always remove on intermediate containers",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_IMAGE_ID_FILE"},
		FilePath: string("/vela/parameters/docker/build/image_id_file,/vela/secrets/docker/build/image_id_file"),
		Name:     "build.image-id-file",
		Usage:    "enables setting writing the image ID to the file",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_ISOLATION"},
		FilePath: string("/vela/parameters/docker/build/isolation,/vela/secrets/docker/build/isolation"),
		Name:     "build.isolation",
		Usage:    "enables container isolation technology",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_LABELS"},
		FilePath: string("/vela/parameters/docker/build/labels,/vela/secrets/docker/build/labels"),
		Name:     "build.labels",
		Usage:    "enables setting metadata for an image",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_MEMORY"},
		FilePath: string("/vela/parameters/docker/build/memory,/vela/secrets/docker/build/memory"),
		Name:     "build.memory",
		Usage:    "enables setting a memory limit",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_MEMORY_SWAPS"},
		FilePath: string("/vela/parameters/docker/build/memory_swaps,/vela/secrets/docker/build/memory_swaps"),
		Name:     "build.memory-swaps",
		Usage:    "enables setting a memory limit",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_NETWORK"},
		FilePath: string("/vela/parameters/docker/build/network,/vela/secrets/docker/build/network"),
		Name:     "build.network",
		Usage:    "enables setting the networking mode for the RUN instructions during build (default \"default\")",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_NO_CACHE"},
		FilePath: string("/vela/parameters/docker/build/no_cache,/vela/secrets/docker/build/no_cache"),
		Name:     "build.no-cache",
		Usage:    "enables setting the networking mode for the RUN instructions during build (default \"default\")",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_OUTPUTS"},
		FilePath: string("/vela/parameters/docker/build/outputs,/vela/secrets/docker/build/outputs"),
		Name:     "build.outputs",
		Usage:    "enables setting an output destination (format: type=local,dest=path)",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_PLATFORM"},
		FilePath: string("/vela/parameters/docker/build/platform,/vela/secrets/docker/build/platform"),
		Name:     "build.platform",
		Usage:    "enables setting a platform if server is multi-platform capable",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_PROGRESS"},
		FilePath: string("/vela/parameters/docker/build/progress,/vela/secrets/docker/build/progress"),
		Name:     "build.progress",
		Usage:    "enables setting type of progress output - options (auto|plain|tty)",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_PULL"},
		FilePath: string("/vela/parameters/docker/build/pull,/vela/secrets/docker/build/pull"),
		Name:     "build.pull",
		Usage:    "enables always attempting to pull a newer version of the image",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_QUIET"},
		FilePath: string("/vela/parameters/docker/build/quiet,/vela/secrets/docker/build/quiet"),
		Name:     "build.quiet",
		Usage:    "enables suppressing the build output and print image ID on success",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_REMOVE"},
		FilePath: string("/vela/parameters/docker/build/remove,/vela/secrets/docker/build/remove"),
		Name:     "build.remove",
		Usage:    "enables removing the intermediate containers after a successful build (default true)",
		Value:    true,
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_SECRETS"},
		FilePath: string("/vela/parameters/docker/build/secrets,/vela/secrets/docker/build/secrets"),
		Name:     "build.secrets",
		Usage:    "enables setting a secret file to expose to the build (only if BuildKit enabled): id=mysecret,src=/local/secret",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_SECURITY_OPTS"},
		FilePath: string("/vela/parameters/docker/build/security_opts,/vela/secrets/docker/build/security_opts"),
		Name:     "build.security-opts",
		Usage:    "enables setting security options",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_SHM_SIZES"},
		FilePath: string("/vela/parameters/docker/build/shm_sizes,/vela/secrets/docker/build/shm_sizes"),
		Name:     "build.shm-sizes",
		Usage:    "enables setting the size of /dev/shm",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_SQUASH"},
		FilePath: string("/vela/parameters/docker/build/squash,/vela/secrets/docker/build/squash"),
		Name:     "build.squash",
		Usage:    "enables setting squash newly built layers into a single new layer",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_SSH_COMPONENTS"},
		FilePath: string("/vela/parameters/docker/build/ssh_components,/vela/secrets/docker/build/ssh_components"),
		Name:     "build.ssh-components",
		Usage:    "enables setting an ssh agent socket or keys to expose to the build (only if BuildKit enabled) (format: default|<id>[=<socket>|<key>[,<key>]])",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"PARAMETER_STREAM"},
		FilePath: string("/vela/parameters/docker/build/stream,/vela/secrets/docker/build/stream"),
		Name:     "build.stream",
		Usage:    "enables streaming attaches to server to negotiate build context",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_TAGS"},
		FilePath: string("/vela/parameters/docker/build/tags,/vela/secrets/docker/build/tags"),
		Name:     "build.tags",
		Usage:    "enables naming and optionally a tag in the 'name:tag' format",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_TARGET"},
		FilePath: string("/vela/parameters/docker/build/target,/vela/secrets/docker/build/target"),
		Name:     "build.target",
		Usage:    "enables setting the target build stage to build.",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_ULIMITS"},
		FilePath: string("/vela/parameters/docker/build/ulimits,/vela/secrets/docker/build/ulimits"),
		Name:     "build.ulimits",
		Usage:    "enables setting ulimit options (default [])",
	},

	// extract vars for open image specification labeling
	&cli.StringFlag{
		EnvVars: []string{"VELA_BUILD_AUTHOR_EMAIL"},
		Name:    "label.author-email",
		Usage:   "author from the source commit",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_BUILD_COMMIT"},
		Name:    "label.commit",
		Usage:   "commit sha from the source commit",
	},
	&cli.IntFlag{
		EnvVars: []string{"VELA_BUILD_NUMBER"},
		Name:    "label.number",
		Usage:   "build number",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_REPO_FULL_NAME"},
		Name:    "label.full-name",
		Usage:   "full name of the repository",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_REPO_LINK"},
		Name:    "label.url",
		Usage:   "direct url of the repository",
	},
}

// Command formats and outputs the Build command from
// the provided configuration to build a Docker image.
// nolint
func (b *Build) Command() (*exec.Cmd, error) {
	logrus.Trace("creating docker build command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if AddHosts is provided
	if len(b.AddHosts) > 0 {
		for _, a := range b.AddHosts {
			// add flag for AddHosts from provided build command
			flags = append(flags, "--add-host", a)
		}
	}

	// check if BuildArgs is provided
	if len(b.BuildArgs) > 0 {
		for _, b := range b.BuildArgs {
			// add flag for BuildArgs from provided build command
			flags = append(flags, "--build-arg", b)
		}
	}

	// check if CacheFrom is provided
	if len(b.CacheFrom) > 0 {
		// add flag for CacheFrom from provided build command
		flags = append(flags, "--cache-from", b.CacheFrom)
	}

	// check if CGroupParent is provided
	if len(b.CGroupParent) > 0 {
		// add flag for CGroupParent from provided build command
		flags = append(flags, "--cgroup-parent", b.CGroupParent)
	}

	// check if Compress is provided
	if b.Compress {
		// add flag for Compress from provided build command
		flags = append(flags, "--compress")
	}

	// add flags for CPU configuration
	flags = append(flags, b.CPU.Flags()...)

	// check if DisableContentTrust is provided
	if b.DisableContentTrust {
		// add flag for DisableContentTrust from provided build command
		flags = append(flags, "--disable-content-trust")
	}

	// check if File is provided
	if len(b.File) > 0 {
		// add flag for File from provided build command
		flags = append(flags, "--file", b.File)
	}

	// check if ForceRM is provided
	if b.ForceRM {
		// add flag for ForceRM from provided build command
		flags = append(flags, "--force-rm")
	}

	// check if ImageIDFile is provided
	if len(b.ImageIDFile) > 0 {
		// add flag for ImageIDFile from provided build command
		flags = append(flags, "--iidfile", b.ImageIDFile)
	}

	// check if Isolation is provided
	if len(b.Isolation) > 0 {
		// add flag for Isolation from provided build command
		flags = append(flags, "--isolation", b.Isolation)
	}

	// check if Labels is provided
	if len(b.Labels) > 0 {
		for _, l := range b.Labels {
			// add flag for Labels from provided build command
			flags = append(flags, "--label", l)
		}
	}

	// check if Memory is provided
	if len(b.Memory) > 0 {
		for _, m := range b.Memory {
			// add flag for Memories from provided build command
			flags = append(flags, "--memory", m)
		}
	}

	// check if MemorySwaps is provided
	if len(b.MemorySwaps) > 0 {
		for _, m := range b.MemorySwaps {
			// add flag for Memories from provided build command
			flags = append(flags, "--memory-swap", m)
		}
	}

	// check if Network is provided
	if len(b.Network) > 0 {
		// add flag for Network from provided build command
		flags = append(flags, "--network", b.Network)
	}

	// check if NoCache is provided
	if b.NoCache {
		// add flag for NoCache from provided build command
		flags = append(flags, "--no-cache")
	}

	// check if Outputs is provided
	if len(b.Outputs) > 0 {
		for _, o := range b.Outputs {
			// add flag for Outputs from provided build command
			flags = append(flags, "--output", o)
		}
	}

	// check if Platform is provided
	if len(b.Platform) > 0 {
		// add flag for Platform from provided build command
		flags = append(flags, "--platform", b.Platform)
	}

	// check if Progress is provided
	if len(b.Progress) > 0 {
		// add flag for Progress from provided build command
		flags = append(flags, "--progress", b.Progress)
	}

	// check if Pull is provided
	if b.Pull {
		// add flag for Pull from provided build command
		flags = append(flags, "--pull")
	}

	// check if Quiet is provided
	if b.Quiet {
		// add flag for Quiet from provided build command
		flags = append(flags, "--quiet")
	}

	// check if Remove is provided
	if b.Remove {
		// add flag for Remove from provided build command
		flags = append(flags, "--rm")
	}

	// check if Secrets is provided
	if len(b.Secrets) > 0 {
		for _, s := range b.Secrets {
			// add flag for Secrets from provided build command
			flags = append(flags, "--secret", s)
		}
	}

	// check if SecurityOpts is provided
	if len(b.SecurityOpts) > 0 {
		for _, s := range b.SecurityOpts {
			// add flag for SecurityOpts from provided build command
			flags = append(flags, "--security-opt", s)
		}
	}

	// check if ShmSizes is provided
	if len(b.ShmSizes) > 0 {
		for _, s := range b.ShmSizes {
			// add flag for ShmSizes from provided build command
			flags = append(flags, "--shm-size", s)
		}
	}

	// check if Squash is provided
	if b.Squash {
		// add flag for Squash from provided build command
		flags = append(flags, "--squash")
	}

	// check if SshComponents is provided
	if len(b.SshComponents) > 0 {
		for _, s := range b.SshComponents {
			// add flag for SshComponents from provided build command
			flags = append(flags, "--ssh", s)
		}
	}

	// check if Stream is provided
	if b.Stream {
		// add flag for Stream from provided build command
		flags = append(flags, "--stream")
	}

	// check if Tags is provided
	if len(b.Tags) > 0 {
		for _, t := range b.Tags {
			// add flag for Tags from provided build command
			flags = append(flags, "--tag", t)
		}
	}

	// check if Target is provided
	if len(b.Target) > 0 {
		// add flag for Target from provided build command
		flags = append(flags, "--target", b.Target)
	}

	// check if Ulimits is provided
	if len(b.Ulimits) > 0 {
		for _, u := range b.Ulimits {
			// add flag for Ulimits from provided build command
			flags = append(flags, "--ulimit", u)
		}
	}

	// add the required directory param
	flags = append(flags, b.Context)

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_docker, append([]string{buildAction}, flags...)...), nil
}

// Exec formats and runs the commands for building a Docker image.
func (b *Build) Exec() error {
	logrus.Trace("running build with provided configuration")

	// add standardized image labels
	b.Labels = append(b.Labels, b.AddLabels()...)

	// create the build command for the file
	cmd, err := b.Command()
	if err != nil {
		return err
	}

	// run the build command for the file
	err = execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// AddLabels adds open container spec labels to plugin
//
// https://github.com/opencontainers/image-spec/blob/v1.0.1/annotations.md
func (b *Build) AddLabels() []string {
	return []string{
		fmt.Sprintf("org.opencontainers.image.created=%s", b.Label.Created),
		fmt.Sprintf("org.opencontainers.image.url=%s", b.Label.URL),
		fmt.Sprintf("org.opencontainers.image.revision=%s", b.Label.Commit),
		fmt.Sprintf("io.vela.build.author=%s", b.Label.AuthorEmail),
		fmt.Sprintf("io.vela.build.number=%d", b.Label.Number),
		fmt.Sprintf("io.vela.build.repo=%s", b.Label.FullName),
		fmt.Sprintf("io.vela.build.commit=%s", b.Label.Commit),
		fmt.Sprintf("io.vela.build.url=%s", b.Label.URL),
	}
}

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (b *Build) Unmarshal() error {
	logrus.Trace("unmarshaling build options")

	// allocate structs to store CPU configuration
	b.CPU = &CPU{}

	// check if any docker options were passed
	if len(b.CPURaw) > 0 {
		// cast raw cpu options into bytes
		cpuOpts := []byte(b.CPURaw)

		// serialize raw cpu options into expected CPU type
		err := json.Unmarshal(cpuOpts, &b.CPU)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build plugin configuration")

	// alert user context is defaulted
	if strings.EqualFold(b.Context, ".") {
		logrus.Warn("running build in default context")
	}

	// verify tag are provided
	if len(b.Tags) == 0 {
		return fmt.Errorf("no build tags provided")
	}

	//TODO Add validation to fields that have custom syntax

	return nil
}

// Flags formats and outputs the flags for
// configuring a Docker.
func (c *CPU) Flags() []string {
	// variable to store flags for command
	var flags []string

	// check if Period is provided
	if c.Period > 0 {
		// add flag for Period from provided build command
		flags = append(flags, "--cpu-period", strconv.Quote(strconv.Itoa(c.Period)))
	}

	// check if Quota is provided
	if c.Quota > 0 {
		// add flag for Quota from provided build command
		flags = append(flags, "--cpu-quota", strconv.Quote(strconv.Itoa(c.Quota)))
	}

	// check if Shares is provided
	if c.Shares > 0 {
		// add flag for Shares from provided build command
		flags = append(flags, "--cpu-shares", strconv.Quote(strconv.Itoa(c.Shares)))
	}

	// check if SetCpus is provided
	if len(c.SetCpus) > 0 {
		// add flag for SetCpus from provided build command
		flags = append(flags, "--cpuset-cpus", strconv.Quote(c.SetCpus))
	}

	// check if SetMems is provided
	if len(c.SetMems) > 0 {
		// add flag for SetMems from provided build command
		flags = append(flags, "--cpuset-mems", strconv.Quote(c.SetMems))
	}

	return flags
}
