package config

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		EnvVars: []string{"DEBUG_MODE", "DEBUG"},
		Name:    "debug",
		Usage:   "enable app debug mode",
		Value:   false,
	},
	&cli.StringFlag{
		EnvVars: []string{"SERVER_ADDR"},
		Name:    "server-addr",
		Aliases: []string{"s"},
		Usage:   "server address",
		Value:   ":8000",
	},
	&cli.StringFlag{
		EnvVars: []string{"CONFIG", "CONFIG_FILE"},
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Load configuration from `FILE`",
		Value:   "/etc/redis-metrics.yaml",
	},
	&cli.IntFlag{
		Name:    "collector-interval",
		Aliases: []string{"i"},
		EnvVars: []string{"COLLECTOR_INTERVAL", "COLLECTOR_INTERVAL_SECONDS"},
		Usage:   "interval seconds of Redis collector scrape.",
		Value:   60,
	},
}
