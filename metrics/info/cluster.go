package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type RedisClusterCollector struct {
	ClusterEnabled *prometheus.GaugeVec
}

func NewRedisClusterCollector() *RedisClusterCollector {
	var (
		// Cluster
		// cluster_enabled:0
		redisClusterClusterEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cluster",
			Name:      "cluster_enabled",
			Help:      "Indicate Redis cluster is enabled.",
		},
			[]string{"node_name", "node_address"})
	)
	return &RedisClusterCollector{
		redisClusterClusterEnabled,
	}
}

func (m *RedisClusterCollector) MustRegister(registry *prometheus.Registry) {
	// Cluster
	registry.MustRegister(m.ClusterEnabled)
}

func (m *RedisClusterCollector) Unregister(registry *prometheus.Registry) bool {
	// Cluster
	if !registry.Unregister(m.ClusterEnabled) {
		return false
	}
	return true
}

func (m *RedisClusterCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	if clusterEnabledStr, ok := r["cluster_enabled"]; ok {
		if clusterEnabled, err := strconv.Atoi(clusterEnabledStr); err == nil {
			m.ClusterEnabled.WithLabelValues(nodeName, nodeAddress).Set(float64(clusterEnabled))
		}
	}
	return nil
}
