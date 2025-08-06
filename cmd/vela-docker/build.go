// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

const buildAction = "build"

type (
	// Build represents the plugin configuration for build information.
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
		Output string
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
		// enables setting the Docker repository name for the image
		Repo string
		// enables setting a secret file to expose to the build (only if BuildKit enabled): id=mysecret,src=/local/secret
		Secret string
		// enables setting security options
		SecurityOpts []string
		// enables setting the size of /dev/shm
		ShmSizes []string
		// enables setting squash newly built layers into a single new layer
		Squash bool
		// enables setting an ssh agent socket or keys to expose to the build (only if BuildKit enabled) (format: default|<id>[=<socket>|<key>[,<key>]])
		SSHComponents []string
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
var buildFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:  "build.add-hosts",
		Usage: "enables adding a custom host-to-IP mapping (host:ip)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_ADD_HOSTS"),
			cli.EnvVar("DOCKER_ADD_HOSTS"),
			cli.File("/vela/parameters/docker/add_hosts"),
			cli.File("/vela/secrets/docker/add_hosts"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.build-args",
		Usage: "enables setting build time arguments for the dockerfile",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_BUILD_ARGS"),
			cli.EnvVar("DOCKER_BUILD_ARGS"),
			cli.File("/vela/parameters/docker/build_args"),
			cli.File("/vela/secrets/docker/build_args"),
		),
	},
	&cli.StringFlag{
		Name:  "build.cache-from",
		Usage: "enables setting images to consider as cache sources",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_CACHE_FROM"),
			cli.EnvVar("DOCKER_CACHE_FROM"),
			cli.File("/vela/parameters/docker/cache_from"),
			cli.File("/vela/secrets/docker/cache_from"),
		),
	},
	&cli.StringFlag{
		Name:  "build.cgroup-parent",
		Usage: "enables setting an optional parent cgroup for the container",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_CGROUP_PARENT"),
			cli.EnvVar("DOCKER_CGROUP_PARENT"),
			cli.File("/vela/parameters/docker/cgroup_parent"),
			cli.File("/vela/secrets/docker/cgroup_parent"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.compress",
		Usage: "enables setting compression the build context using gzip",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_COMPRESS"),
			cli.EnvVar("DOCKER_COMPRESS"),
			cli.File("/vela/parameters/docker/compress"),
			cli.File("/vela/secrets/docker/compress"),
		),
	},
	&cli.StringFlag{
		Name:  "build.context",
		Value: ".",
		Usage: "enables setting the build context",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_CONTEXT"),
			cli.EnvVar("DOCKER_CONTEXT"),
			cli.File("/vela/parameters/docker/context"),
			cli.File("/vela/secrets/docker/context"),
		),
	},
	&cli.StringFlag{
		Name:  "build.cpu",
		Usage: "enables setting custom cpu options",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_CPU"),
			cli.EnvVar("DOCKER_CPU"),
			cli.File("/vela/parameters/docker/cpu"),
			cli.File("/vela/secrets/docker/cpu"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.disable-content-trust",
		Usage: "enables skipping image verification (default true)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_DISABLE_CONTENT_TRUST"),
			cli.EnvVar("DOCKER_DISABLE_CONTENT_TRUST"),
			cli.File("/vela/parameters/docker/disable-content-trust"),
			cli.File("/vela/secrets/docker/disable-content-trust"),
		),
	},
	&cli.StringFlag{
		Name:  "build.file",
		Usage: "enables setting the name of the Dockerfile (Default is 'PATH/Dockerfile')",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_FILE"),
			cli.EnvVar("DOCKER_FILE"),
			cli.File("/vela/parameters/docker/file"),
			cli.File("/vela/secrets/docker/file"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.force-rm",
		Usage: "enables setting always remove on intermediate containers",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_FORCE_RM"),
			cli.EnvVar("DOCKER_FORCE_RM"),
			cli.File("/vela/parameters/docker/force_rm"),
			cli.File("/vela/secrets/docker/force_rm"),
		),
	},
	&cli.StringFlag{
		Name:  "build.image-id-file",
		Usage: "enables setting writing the image ID to the file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_IMAGE_ID_FILE"),
			cli.EnvVar("DOCKER_IMAGE_ID_FILE"),
			cli.File("/vela/parameters/docker/image_id_file"),
			cli.File("/vela/secrets/docker/image_id_file"),
		),
	},
	&cli.StringFlag{
		Name:  "build.isolation",
		Usage: "enables container isolation technology",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_ISOLATION"),
			cli.EnvVar("DOCKER_ISOLATION"),
			cli.File("/vela/parameters/docker/isolation"),
			cli.File("/vela/secrets/docker/isolation"),
		),
	},
	&cli.StringFlag{
		Name:  "build.labels",
		Usage: "enables setting metadata for an image",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_LABELS"),
			cli.EnvVar("DOCKER_LABELS"),
			cli.File("/vela/parameters/docker/labels"),
			cli.File("/vela/secrets/docker/labels"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.memory",
		Usage: "enables setting a memory limit",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_MEMORY"),
			cli.EnvVar("DOCKER_MEMORY"),
			cli.File("/vela/parameters/docker/memory"),
			cli.File("/vela/secrets/docker/memory"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.memory-swaps",
		Usage: "enables setting a memory limit",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_MEMORY_SWAPS"),
			cli.EnvVar("DOCKER_MEMORY_SWAPS"),
			cli.File("/vela/parameters/docker/memory_swaps"),
			cli.File("/vela/secrets/docker/memory_swaps"),
		),
	},
	&cli.StringFlag{
		Name:  "build.network",
		Usage: "enables setting the networking mode for the RUN instructions during build (default \"default\")",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_NETWORK"),
			cli.EnvVar("DOCKER_NETWORK"),
			cli.File("/vela/parameters/docker/network"),
			cli.File("/vela/secrets/docker/network"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.no-cache",
		Usage: "enables setting the networking mode for the RUN instructions during build (default \"default\")",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_NO_CACHE"),
			cli.EnvVar("DOCKER_NO_CACHE"),
			cli.File("/vela/parameters/docker/no_cache"),
			cli.File("/vela/secrets/docker/no_cache"),
		),
	},
	&cli.StringFlag{
		Name:  "build.output",
		Usage: "set an output destination (format: type=local,dest=path)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_OUTPUT"),
			cli.EnvVar("DOCKER_OUTPUT"),
			cli.File("/vela/parameters/docker/output"),
			cli.File("/vela/secrets/docker/output"),
		),
	},
	&cli.StringFlag{
		Name:  "build.platform",
		Usage: "enables setting a platform if server is multi-platform capable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_PLATFORM"),
			cli.EnvVar("DOCKER_PLATFORM"),
			cli.File("/vela/parameters/docker/platform"),
			cli.File("/vela/secrets/docker/platform"),
		),
	},
	&cli.StringFlag{
		Name:  "build.progress",
		Usage: "enables setting type of progress output - options (auto|plain|tty)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_PROGRESS"),
			cli.EnvVar("DOCKER_PROGRESS"),
			cli.File("/vela/parameters/docker/progress"),
			cli.File("/vela/secrets/docker/progress"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.pull",
		Usage: "enables always attempting to pull a newer version of the image",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_PULL"),
			cli.EnvVar("DOCKER_PULL"),
			cli.File("/vela/parameters/docker/pull"),
			cli.File("/vela/secrets/docker/pull"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.quiet",
		Usage: "enables suppressing the build output and print image ID on success",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_QUIET"),
			cli.EnvVar("DOCKER_QUIET"),
			cli.File("/vela/parameters/docker/quiet"),
			cli.File("/vela/secrets/docker/quiet"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.remove",
		Value: true,
		Usage: "enables removing the intermediate containers after a successful build (default true)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_REMOVE"),
			cli.EnvVar("DOCKER_REMOVE"),
			cli.File("/vela/parameters/docker/remove"),
			cli.File("/vela/secrets/docker/remove"),
		),
	},
	&cli.StringFlag{
		Name:  "build.repo",
		Usage: "Docker repository name for the image",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_REPO"),
			cli.EnvVar("DOCKER_REPO"),
			cli.File("/vela/parameters/docker/repo"),
			cli.File("/vela/secrets/docker/repo"),
		),
	},
	&cli.StringFlag{
		Name:  "build.secret",
		Usage: "set a secret file to expose to the build (only if BuildKit enabled): id=mysecret,src=/local/secret",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_SECRET"),
			cli.EnvVar("DOCKER_SECRET"),
			cli.File("/vela/parameters/docker/secret"),
			cli.File("/vela/secrets/docker/secret"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.security-opts",
		Usage: "enables setting security options",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_SECURITY_OPTS"),
			cli.EnvVar("DOCKER_SECURITY_OPTS"),
			cli.File("/vela/parameters/docker/security_opts"),
			cli.File("/vela/secrets/docker/security_opts"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.shm-sizes",
		Usage: "enables setting the size of /dev/shm",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_SHM_SIZES"),
			cli.EnvVar("DOCKER_SHM_SIZES"),
			cli.File("/vela/parameters/docker/shm_sizes"),
			cli.File("/vela/secrets/docker/shm_sizes"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.squash",
		Usage: "enables setting squash newly built layers into a single new layer",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_SQUASH"),
			cli.EnvVar("DOCKER_SQUASH"),
			cli.File("/vela/parameters/docker/squash"),
			cli.File("/vela/secrets/docker/squash"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.ssh-components",
		Usage: "enables setting an ssh agent socket or keys to expose to the build (only if BuildKit enabled) (format: default|<id>[=<socket>|<key>[,<key>]])",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_SSH_COMPONENTS"),
			cli.EnvVar("DOCKER_SSH_COMPONENTS"),
			cli.File("/vela/parameters/docker/ssh_components"),
			cli.File("/vela/secrets/docker/ssh_components"),
		),
	},
	&cli.BoolFlag{
		Name:  "build.stream",
		Usage: "enables streaming attaches to server to negotiate build context",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_STREAM"),
			cli.EnvVar("DOCKER_STREAM"),
			cli.File("/vela/parameters/docker/stream"),
			cli.File("/vela/secrets/docker/stream"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.tags",
		Usage: "enables naming and optionally a tag in the 'name:tag' format",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_TAGS"),
			cli.EnvVar("DOCKER_TAGS"),
			cli.File("/vela/parameters/docker/tags"),
			cli.File("/vela/secrets/docker/tags"),
		),
	},
	&cli.StringFlag{
		Name:  "build.target",
		Usage: "enables setting the target build stage to build.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_TARGET"),
			cli.EnvVar("DOCKER_TARGET"),
			cli.File("/vela/parameters/docker/target"),
			cli.File("/vela/secrets/docker/target"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "build.ulimits",
		Usage: "enables setting ulimit options (default [])",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PARAMETER_ULIMITS"),
			cli.EnvVar("DOCKER_ULIMITS"),
			cli.File("/vela/parameters/docker/ulimits"),
			cli.File("/vela/secrets/docker/ulimits"),
		),
	},

	// extract vars for open image specification labeling
	&cli.StringFlag{
		Name:    "label.author-email",
		Usage:   "author from the source commit",
		Sources: cli.EnvVars("VELA_BUILD_AUTHOR_EMAIL"),
	},
	&cli.StringFlag{
		Name:    "label.commit",
		Usage:   "commit sha from the source commit",
		Sources: cli.EnvVars("VELA_BUILD_COMMIT"),
	},
	&cli.IntFlag{
		Name:    "label.number",
		Usage:   "build number",
		Sources: cli.EnvVars("VELA_BUILD_NUMBER"),
	},
	&cli.StringFlag{
		Name:    "label.full-name",
		Usage:   "full name of the repository",
		Sources: cli.EnvVars("VELA_REPO_FULL_NAME"),
	},
	&cli.StringFlag{
		Name:    "label.url",
		Usage:   "direct url of the repository",
		Sources: cli.EnvVars("VELA_REPO_LINK"),
	},
}

// Command formats and outputs the Build command from
// the provided configuration to build a Docker image.
//
//nolint:gocyclo // Ignore line length
func (b *Build) Command(ctx context.Context) *exec.Cmd {
	logrus.Trace("creating docker build command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// iterate through the additional hosts provided
	for _, a := range b.AddHosts {
		// add flag for AddHosts from provided build command
		flags = append(flags, "--add-host", a)
	}

	// iterate through the build arguments provided
	for _, b := range b.BuildArgs {
		// add flag for BuildArgs from provided build command
		flags = append(flags, "--build-arg", b)
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

	// iterate through the labels provided
	for _, l := range b.Labels {
		// add flag for Labels from provided build command
		flags = append(flags, "--label", l)
	}

	// iterate through the memory arguments provided
	for _, m := range b.Memory {
		// add flag for Memory from provided build command
		flags = append(flags, "--memory", m)
	}

	// iterate through the memory swap arguments provided
	for _, m := range b.MemorySwaps {
		// add flag for Memory Swaps from provided build command
		flags = append(flags, "--memory-swap", m)
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

	// check if Output is provided
	if len(b.Output) > 0 {
		// add flag for output from provided build command
		flags = append(flags, "--output", b.Output)
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

	// check if Secret is provided
	if len(b.Secret) > 0 {
		// add flag for secret from provided build command
		flags = append(flags, "--secret", b.Secret)
	}

	// iterate through the security options provided
	for _, s := range b.SecurityOpts {
		// add flag for SecurityOpts from provided build command
		flags = append(flags, "--security-opt", s)
	}

	// iterate through the SHM sizes provided
	for _, s := range b.ShmSizes {
		// add flag for ShmSizes from provided build command
		flags = append(flags, "--shm-size", s)
	}

	// check if Squash is provided
	if b.Squash {
		// add flag for Squash from provided build command
		flags = append(flags, "--squash")
	}

	// iterate through the SSH components provided
	for _, s := range b.SSHComponents {
		// add flag for SSHComponents from provided build command
		flags = append(flags, "--ssh", s)
	}

	// check if Stream is provided
	if b.Stream {
		// add flag for Stream from provided build command
		flags = append(flags, "--stream")
	}

	// iterate through the tags provided
	for _, t := range b.Tags {
		// check if a Docker repository was provided
		if len(b.Repo) > 0 {
			// check if the tag already has the repo in it
			if !strings.Contains(t, b.Repo) {
				t = fmt.Sprintf("%s:%s", b.Repo, t)
			}
		}

		// add flag for Tags from provided build command
		flags = append(flags, "--tag", t)
	}

	// check if Target is provided
	if len(b.Target) > 0 {
		// add flag for Target from provided build command
		flags = append(flags, "--target", b.Target)
	}

	// iterate through the ulimits provided
	for _, u := range b.Ulimits {
		// add flag for Ulimits from provided build command
		flags = append(flags, "--ulimit", u)
	}

	// add the required directory param
	flags = append(flags, b.Context)

	//nolint:gosec // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.CommandContext(ctx, _docker, append([]string{buildAction}, flags...)...)
}

// Exec formats and runs the commands for building a Docker image.
func (b *Build) Exec(ctx context.Context) error {
	logrus.Trace("running build with provided configuration")

	// add standardized image labels
	b.Labels = append(b.Labels, b.AddLabels()...)

	// create the build command for the file
	cmd := b.Command(ctx)

	// run the build command for the file
	err := execCmd(cmd)
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
