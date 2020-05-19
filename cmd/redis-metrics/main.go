package main

import (
	"os"

	"github.com/cwr0401/redis_metrics/config"
	"github.com/cwr0401/redis_metrics/metrics"
	"github.com/cwr0401/redis_metrics/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "redis-metrics",
		Version: version.Version.String(),
		Flags:   config.Flags,
		Action:  metrics.RedisMetricsAction,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Chen Weiran",
				Email: "cwr0401@gmail.com",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
