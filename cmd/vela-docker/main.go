// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-vela/vela-docker/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-docker"
	app.HelpName = "vela-docker"
	app.Usage = "Vela Docker plugin for building and publishing images"
	app.Copyright = "Copyright (c) 2022 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "VELA_LOG_LEVEL", "DOCKER_LOG_LEVEL"},
			FilePath: string("/vela/parameters/docker/log_level,/vela/secrets/docker/log_level"),
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},
	}

	// add build flags
	app.Flags = append(app.Flags, buildFlags...)

	// add daemon flags
	app.Flags = append(app.Flags, daemonFlags...)

	// add push flags
	app.Flags = append(app.Flags, pushFlags...)

	// add registry flags
	app.Flags = append(app.Flags, registryFlags...)

	// add manifest flags
	app.Flags = append(app.Flags, manifestFlags...)

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-docker",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/docker",
		"registry": "https://hub.docker.com/r/target/vela-docker",
	}).Info("Vela Docker Plugin")

	// create the plugin
	p := Plugin{
		Build: &Build{
			AddHosts:            c.StringSlice("build.add-hosts"),
			BuildArgs:           c.StringSlice("build.build-args"),
			CacheFrom:           c.String("build.cache-from"),
			CGroupParent:        c.String("build.cgroup-parent"),
			Compress:            c.Bool("build.compress"),
			Context:             c.String("build.context"),
			CPURaw:              c.String("build.cpu"),
			DisableContentTrust: c.Bool("build.disable-content-trust"),
			File:                c.String("build.file"),
			ForceRM:             c.Bool("build.force-rm"),
			ImageIDFile:         c.String("build.image-id-file"),
			Isolation:           c.String("build.isolation"),
			Label: &Label{
				AuthorEmail: c.String("label.author-email"),
				Commit:      c.String("label.commit"),
				Created:     time.Now().Format(time.RFC3339),
				FullName:    c.String("label.full-name"),
				Number:      c.Int("label.number"),
				URL:         c.String("label.url"),
			},
			Labels:        c.StringSlice("build.labels"),
			Memory:        c.StringSlice("build.memory"),
			MemorySwaps:   c.StringSlice("build.memory-swaps"),
			Network:       c.String("build.network"),
			NoCache:       c.Bool("build.no-cache"),
			Output:        c.String("build.output"),
			Platform:      c.String("build.platform"),
			Progress:      c.String("build.progress"),
			Pull:          c.Bool("build.pull"),
			Quiet:         c.Bool("build.quiet"),
			Remove:        c.Bool("build.remove"),
			Repo:          c.String("build.repo"),
			Secret:        c.String("build.secret"),
			SecurityOpts:  c.StringSlice("build.security-opts"),
			ShmSizes:      c.StringSlice("build.shm-sizes"),
			Squash:        c.Bool("build.squash"),
			SshComponents: c.StringSlice("build.ssh-components"),
			Stream:        c.Bool("build.stream"),
			Tags:          c.StringSlice("build.tags"),
			Target:        c.String("build.target"),
			Ulimits:       c.StringSlice("build.ulimits"),
		},
		Daemon: &Daemon{},
		Push: &Push{
			DisableContentTrust: c.Bool("push.disable-content-trust"),
		},
		Registry: &Registry{
			DryRun:   c.Bool("registry.dry-run"),
			Name:     c.String("registry.name"),
			Password: c.String("registry.password"),
			Username: c.String("registry.username"),
		},
		Manifest: &Manifest{
			RawSpec: c.String("manifest.spec"),
		},
	}

	// validate the plugin
	err := p.Validate(c.String("daemon"))
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
