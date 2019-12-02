package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var (
	// Cluster
	// cluster_enabled:0
	RedisClusterClusterEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "cluster",
		Name:      "cluster_enabled",
		Help:      "Indicate Redis cluster is enabled.",
	},
		[]string{"node_name", "node_address"})
)

func SetRedisCluster(nodeName, nodeAddress string, r map[string]string) error {
	if clusterEnabledStr, ok := r["cluster_enabled"]; ok {
		if clusterEnabled, err := strconv.Atoi(clusterEnabledStr); err == nil {
			RedisClusterClusterEnabled.WithLabelValues(nodeName, nodeAddress).Set(float64(clusterEnabled))
		}
	}
	return nil
}
