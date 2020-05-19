package info

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	RedisInfoSections = []string{
		"Server", "Clients", "Memory", "Persistence", "Stats",
		"Replication", "CPU", "Commandstats", "Cluster", "Keyspace",
	}
)

type Collector interface {
	MustRegister(registry *prometheus.Registry)
	Unregister(registry *prometheus.Registry) bool
	Set(nodeName, nodeAddress string, r map[string]string) error
}

type RedisCollector map[string]Collector

func NewRedisCollector() RedisCollector {
	return RedisCollector{
		"Server":       NewRedisServerCollector(),
		"Clients":      NewRedisClientsCollector(),
		"Memory":       NewRedisMemoryCollector(),
		"Persistence":  NewRedisPersistenceCollector(),
		"Stats":        NewRedisStatsCollector(),
		"Replication":  NewRedisReplicationCollector(),
		"CPU":          NewRedisCPUCollector(),
		"Commandstats": NewRedisCommandstatsCollector(),
		"Cluster":      NewRedisClusterCollector(),
		"Keyspace":     NewRedisKeyspaceCollector(),
		"Sentinel":     NewRedisSentinelCollector(),
	}
}

func (m RedisCollector) MustRegister(registry *prometheus.Registry) {
	// log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	// log.SetLevel(log.DebugLevel)
	//var (
	//	registry                        = prometheus.NewRegistry()
	//	gather      prometheus.Gatherer = registry
	//	promHandler                     = promhttp.HandlerFor(gather, promhttp.HandlerOpts{})
	//	Handler                         = promhttp.InstrumentMetricHandler(registry, promHandler)
	//)

	// Collector have to be registered to be exposed:
	// registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	// registry.MustRegister(prometheus.NewGoCollector())
	for key, metrics := range m {
		log.Infof("Register %s", key)
		metrics.MustRegister(registry)
	}
}

func (m RedisCollector) Unregister(registry *prometheus.Registry) bool {
	for key, metrics := range m {
		log.Infof("Unregister %s", key)
		if !metrics.Unregister(registry) {
			return false
		}
	}
	return true
}

func (m RedisCollector) Set(nodeName, nodeAddress string, r map[string]string) error {

	for _, section := range RedisInfoSections {
		log.WithFields(log.Fields{
			"node": nodeName,
			"addr": nodeAddress,
		}).Infof("Set metrics %s", section)
		if metrics, ok := m[section]; ok {
			err := metrics.Set(nodeName, nodeAddress, r)
			if err != nil {
				return err
			}
		} else {
			log.WithFields(log.Fields{
				"node": nodeName,
				"addr": nodeAddress,
			}).Errorf("The metrics %s not found", section)
			continue
		}
	}
	return nil
}
