package metrics

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cwr0401/redis_metrics/config"
	"github.com/cwr0401/redis_metrics/metrics/info"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	MinCollectorInterval int = 20
	MaxCollectorInterval int = 600
)

var reloadChan = make(chan struct{}, 1)

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

func reload(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" && r.Method != "PUT" {
		w.WriteHeader(405)
		w.Write([]byte("Only POST or PUT requests allowed"))
		return
	}

	if strings.Contains(r.RemoteAddr, "127.0.0.1:") {
		go func() {
			reloadChan <- struct{}{}
		}()
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(405)
		w.Write([]byte("Error"))
	}
}

func HttpServer(addr string, handler http.Handler) {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/-/reload", reload)
	http.Handle("/metrics", handler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panic(err)
	}
}

func readConfigFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func loadConfigFile(path string) (*config.RedisConfig, error) {

	configFilePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	log.Infof("Loading redis-metrics configuration file %s", configFilePath)
	configContent, err := readConfigFile(configFilePath)
	if err != nil {
		return nil, err
	}
	redisConfig, err := config.ParseRedisConfig(configContent)
	if err != nil {
		return nil, err
	}

	return redisConfig, nil
}

func reloadConfig(configFile string, redisConfig *config.RedisConfig) *config.RedisConfig {
	for {
		<-reloadChan
		log.Info("Reload configuration file")
		newRedisConfig, err := loadConfigFile(configFile)
		if err != nil {
			log.Errorf("Reload configuration file error: %s", err)
			continue
		}
		if newRedisConfig.Equal(*redisConfig) {
			log.Info("config file no changes")
			continue
		} else {
			return newRedisConfig
		}
	}
}

func RedisMetricsAction(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Set logger Debug mode")
	}
	addr := c.String("server-addr")
	log.Infof("Http listen on %s", addr)

	collectorInterval := c.Int("collector-interval")
	if collectorInterval < MinCollectorInterval {
		collectorInterval = MinCollectorInterval
	}
	if collectorInterval > MaxCollectorInterval {
		collectorInterval = MaxCollectorInterval
	}
	collectorIntervalDuration := time.Duration(collectorInterval) * time.Second
	log.Infof("Scrape interval duration %s", collectorIntervalDuration)

	configFile := c.String("config")
	redisConfig, err := loadConfigFile(configFile)
	if err != nil {
		return err
	}

	var (
		registry                        = prometheus.NewRegistry()
		gather      prometheus.Gatherer = registry
		promHandler                     = promhttp.HandlerFor(gather, promhttp.HandlerOpts{})
		Handler                         = promhttp.InstrumentMetricHandler(registry, promHandler)
	)
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	go HttpServer(addr, Handler)

	for {
		ctx, cancel := context.WithCancel(context.Background())

		// register
		rmc := info.NewRedisCollector()
		rmc.MustRegister(registry)

		rm := RedisMetrics{
			Collector: rmc,
			Duration:  collectorIntervalDuration,
			Config:    redisConfig,
		}

		stopFlag := make(chan struct{}, 1)

		log.Debug("Start Collector")
		go rm.Run(ctx, stopFlag)

		// wait reload event
		redisConfig = reloadConfig(configFile, redisConfig)
		cancel()
		log.Debug("Stop Collector")
		<-stopFlag
		if !rmc.Unregister(registry) {
			return errors.New("unregister error")
		}
	}
}
