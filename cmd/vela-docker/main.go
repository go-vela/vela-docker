// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-docker"
	app.HelpName = "vela-docker"
	app.Usage = "Vela img plugin for building and publishing images"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Compiled = time.Now()
	app.Action = run

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

	err := app.Run(os.Args)
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

	return nil
}
