package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

// cmdstat_client:calls=8,usec=32,usec_per_call=4.00
type RedisCommandstatsCollector struct {
	Calls       *prometheus.GaugeVec
	Usec        *prometheus.GaugeVec
	UsecPerCall *prometheus.GaugeVec
}

func NewRedisCommandstatsCollector() *RedisCommandstatsCollector {
	var (
		redisCommandstatsCalls = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cmdstat",
			Name:      "calls",
			Help:      "",
		}, []string{"node_name", "node_address", "cmd"})

		redisCommandstatsUsec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cmdstat",
			Name:      "usec",
			Help:      "",
		}, []string{"node_name", "node_address", "cmd"})

		redisCommandstatsUsecPerCall = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cmdstat",
			Name:      "usec_per_call",
			Help:      "",
		}, []string{"node_name", "node_address", "cmd"})
	)
	return &RedisCommandstatsCollector{
		redisCommandstatsCalls,
		redisCommandstatsUsec,
		redisCommandstatsUsecPerCall,
	}
}

func (m *RedisCommandstatsCollector) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(m.Calls)
	registry.MustRegister(m.Usec)
	registry.MustRegister(m.UsecPerCall)
}

func (m *RedisCommandstatsCollector) Unregister(registry *prometheus.Registry) bool {
	if !registry.Unregister(m.Calls) {
		return false
	}
	if !registry.Unregister(m.Usec) {
		return false
	}
	if !registry.Unregister(m.UsecPerCall) {
		return false
	}
	return true
}

func (m *RedisCommandstatsCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	for key, value := range r {
		if strings.Contains(key, "cmdstat_") {
			valueMap := stringToMap(value)
			if callsStr, ok := valueMap["calls"]; ok {
				if calls, err := strconv.Atoi(callsStr); err == nil {
					m.Calls.WithLabelValues(nodeName, nodeAddress, key).Set(float64(calls))
				}
			}
			if usecStr, ok := valueMap["usec"]; ok {
				if usec, err := strconv.Atoi(usecStr); err == nil {
					m.Usec.WithLabelValues(nodeName, nodeAddress, key).Set(float64(usec))
				}
			}
			if usecPerCallStr, ok := valueMap["usec_per_call"]; ok {
				if usecPerCall, err := strconv.ParseFloat(usecPerCallStr, 64); err == nil {
					m.UsecPerCall.WithLabelValues(nodeName, nodeAddress, key).Set(usecPerCall)
				}
			}
		}
	}
	return nil
}
