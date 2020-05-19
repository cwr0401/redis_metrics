package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

type RedisKeyspaceCollector struct {
	DBKeys    *prometheus.GaugeVec
	DBExpires *prometheus.GaugeVec
	DBAvgTTL  *prometheus.GaugeVec
}

func NewRedisKeyspaceCollector() *RedisKeyspaceCollector {
	var (
		// db0:keys=5266,expires=5213,avg_ttl=1345519
		// db1:keys=1146347,expires=48039,avg_ttl=56346297
		// db2:keys=1720,expires=1178,avg_ttl=3242494
		redisKeyspaceDBKeys = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "keyspace",
			Name:      "db_keys",
			Help:      "Number of keys for each database.",
		},
			[]string{"node_name", "node_address", "database"})

		redisKeyspaceDBExpires = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "keyspace",
			Name:      "db_expires",
			Help:      "Number of expire keys for each database.",
		},
			[]string{"node_name", "node_address", "database"})

		redisKeyspaceDBAvgTTL = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "keyspace",
			Name:      "db_avg_ttl",
			Help:      "Average of ttl for each database.",
		},
			[]string{"node_name", "node_address", "database"})
	)
	return &RedisKeyspaceCollector{
		redisKeyspaceDBKeys,
		redisKeyspaceDBExpires,
		redisKeyspaceDBAvgTTL,
	}
}

func (m *RedisKeyspaceCollector) MustRegister(registry *prometheus.Registry) {
	// Keyspace
	registry.MustRegister(m.DBKeys)
	registry.MustRegister(m.DBExpires)
	registry.MustRegister(m.DBAvgTTL)
}

func (m *RedisKeyspaceCollector) Unregister(registry *prometheus.Registry) bool {
	// Keyspace
	if !registry.Unregister(m.DBKeys) {
		return false
	}
	if !registry.Unregister(m.DBExpires) {
		return false
	}
	if !registry.Unregister(m.DBAvgTTL) {
		return false
	}
	return true
}

func (m *RedisKeyspaceCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	for key, value := range r {
		ok1 := strings.Contains(key, "db")
		ok2 := strings.Contains(value, "keys")
		if ok1 && ok2 {
			valueMap := stringToMap(value)
			if keysStr, ok := valueMap["keys"]; ok {
				if keys, err := strconv.Atoi(keysStr); err == nil {
					m.DBKeys.WithLabelValues(nodeName, nodeAddress, key).Set(float64(keys))
				}
			}
			if expiresStr, ok := valueMap["expires"]; ok {
				if expires, err := strconv.Atoi(expiresStr); err == nil {
					m.DBExpires.WithLabelValues(nodeName, nodeAddress, key).Set(float64(expires))
				}
			}
			if avgTTLStr, ok := valueMap["avg_ttl"]; ok {
				if avgTTL, err := strconv.Atoi(avgTTLStr); err == nil {
					m.DBAvgTTL.WithLabelValues(nodeName, nodeAddress, key).Set(float64(avgTTL))
				}
			}
		}
	}
	return nil
}
