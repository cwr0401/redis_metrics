package metrics

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cwr0401/redis_metrics/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	MinCollectorInterval int = 20
	MaxCollectorInterval int = 600
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func RedisMetricsServer(c *cli.Context) {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/health", healthCheck)
	http.Handle("/metrics", Handler)

	log.Infof("Http listen on %s", c.String("server-addr"))
	err := http.ListenAndServe(c.String("server-addr"), nil)
	if err != nil {
		log.Panic(err)
	}
}

func RedisMetricsAction(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	go RedisMetricsServer(c)

	configFilePath, err := filepath.Abs(c.String("config"))
	if err != nil {
		log.Panic(err)
	}

	log.Infof("The redis-metrics configuration file %s", configFilePath)
	configContent, err := readConfigFile(configFilePath)
	if err != nil {
		panic(err)
	}
	RedisConfig, err := config.ParseRedisConfig(configContent)
	if err != nil {
		panic(err)
	}

	collectorInterval := c.Int("collector-interval")
	if collectorInterval < MinCollectorInterval {
		collectorInterval = MinCollectorInterval
	}
	if collectorInterval > MaxCollectorInterval {
		collectorInterval = MaxCollectorInterval
	}
	collectorIntervalDuration := time.Duration(collectorInterval) * time.Second
	log.Infof("Scrape interval duration %s", collectorIntervalDuration)

	RedisCollector(RedisConfig, collectorIntervalDuration)
	return nil
}

func readConfigFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
