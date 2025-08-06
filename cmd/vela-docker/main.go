// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/vela-docker/version"
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
	app := &cli.Command{
		Name:      "vela-docker",
		Usage:     "Vela Docker plugin for building and publishing images",
		Copyright: "Copyright 2020 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		// Plugin Metadata
		Version: v.Semantic(),
		Action:  run,
	}

	// Plugin Flags
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "log.level",
			Value: "info",
			Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG_LEVEL"),
				cli.EnvVar("VELA_LOG_LEVEL"),
				cli.EnvVar("DOCKER_LOG_LEVEL"),
				cli.File("/vela/parameters/docker/log_level"),
				cli.File("/vela/secrets/docker/log_level"),
			),
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

	err = app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(ctx context.Context, c *cli.Command) error {
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
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/docker",
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
			SSHComponents: c.StringSlice("build.ssh-components"),
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
	}

	// validate the plugin
	err := p.Validate(c.String("daemon"))
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec(ctx)
}
