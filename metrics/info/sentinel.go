package info

import (
	"github.com/prometheus/client_golang/prometheus"
)

type RedisSentinelCollector struct {
	Up                           *prometheus.GaugeVec
	SentinelMasters              *prometheus.GaugeVec
	SentinelTilt                 *prometheus.GaugeVec
	SentinelRunningScripts       *prometheus.GaugeVec
	SentinelScriptsQueueLength   *prometheus.GaugeVec
	SentinelSimulateFailureFlags *prometheus.GaugeVec
	Master                       *prometheus.GaugeVec
	MastersInfo                  *prometheus.CounterVec
	SlavesInfo                   *prometheus.GaugeVec
	SentinelsInfo                *prometheus.GaugeVec
}

func NewRedisSentinelCollector() *RedisSentinelCollector {

	var (
		// Sentinel
		redisSentinelUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "up",
			Help:      "Value is 1 if Redis server alive, 0 otherwise.",
		},
			[]string{"node_name", "node_address"})

		// sentinel_masters:1
		redisSentinelSentinelMasters = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "sentinel_masters",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		//sentinel_tilt:0
		redisSentinelSentinelTilt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "sentinel_tilt",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		//sentinel_running_scripts:0
		redisSentinelSentinelRunningScripts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "sentinel_running_scripts",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		//sentinel_scripts_queue_length:0
		redisSentinelSentinelScriptsQueueLength = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "sentinel_scripts_queue_length",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		//sentinel_simulate_failure_flags:0
		redisSentinelSentinelSimulateFailureFlags = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "sentinel_simulate_failure_flags",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		//master0:name=mymaster,status=ok,address=172.25.1.3:6379,slaves=2,sentinels=3
		redisSentinelMaster = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "master",
			Help:      "",
		},
			[]string{"node_name", "node_address", "index", "name", "status", "address", "slaves", "sentinels"})

		// sentinel master
		redisSentinelMastersInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "masters_info",
			Help:      "",
		},
			[]string{"name", "ip", "port", "runid", "flags"})

		// "link_pending_commands", "link_refcount", "last_ping_sent",
		// "last_ok_ping_reply", "last_ping_reply", "down_after_milliseconds", "info_refresh", "role_reported",
		// "role_reported_time", "config_epoch",  "num_slaves", "num_other_sentinels", "quorum", "failover_timeout",
		// "parallel_syncs"

		redisSentinelSlavesInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "slaves_info",
			Help:      "",
		},
			[]string{})

		redisSentinelSentinelsInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "sentinel",
			Name:      "sentinels_info",
			Help:      "",
		},
			[]string{})
	)
	return &RedisSentinelCollector{
		redisSentinelUp,
		redisSentinelSentinelMasters,
		redisSentinelSentinelTilt,
		redisSentinelSentinelRunningScripts,
		redisSentinelSentinelScriptsQueueLength,
		redisSentinelSentinelSimulateFailureFlags,
		redisSentinelMaster,
		redisSentinelMastersInfo,
		redisSentinelSlavesInfo,
		redisSentinelSentinelsInfo,
	}
}

func (m *RedisSentinelCollector) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(m.Up)
	registry.MustRegister(m.SentinelMasters)
	registry.MustRegister(m.SentinelTilt)
	registry.MustRegister(m.SentinelRunningScripts)
	registry.MustRegister(m.SentinelScriptsQueueLength)
	registry.MustRegister(m.SentinelSimulateFailureFlags)
	registry.MustRegister(m.Master)
	registry.MustRegister(m.MastersInfo)
	registry.MustRegister(m.SlavesInfo)
	registry.MustRegister(m.SentinelsInfo)
}

func (m *RedisSentinelCollector) Unregister(registry *prometheus.Registry) bool {
	if !registry.Unregister(m.Up) {
		return false
	}
	if !registry.Unregister(m.SentinelMasters) {
		return false
	}
	if !registry.Unregister(m.SentinelTilt) {
		return false
	}
	if !registry.Unregister(m.SentinelRunningScripts) {
		return false
	}
	if !registry.Unregister(m.SentinelScriptsQueueLength) {
		return false
	}
	if !registry.Unregister(m.SentinelSimulateFailureFlags) {
		return false
	}
	if !registry.Unregister(m.Master) {
		return false
	}
	if !registry.Unregister(m.MastersInfo) {
		return false
	}
	if !registry.Unregister(m.SlavesInfo) {
		return false
	}
	if !registry.Unregister(m.SentinelsInfo) {
		return false
	}
	return true
}

func (m *RedisSentinelCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	return nil
}

func (m *RedisSentinelCollector) SetRedisSentinelMastersInfo(r map[string]string) error {
	m.MastersInfo.WithLabelValues(r["name"], r["ip"], r["port"], r["runid"], r["flags"]).Inc()

	// r["link-pending-commands"], r["link-refcount"],
	//		r["last-ping-sent"], r["last-ok-ping-reply"], r["last-ping-reply"], r["down-after-milliseconds"],
	//		r["info-refresh"], r["role-reported"], r["role-reported-time"], r["config-epoch"],  r["num-slaves"],
	//		r["num-other-sentinels"], r["quorum"], r["failover-timeout"], r["parallel-syncs"],

	return nil
}

func SetRedisSentinelSlavesInfo(r map[string]string) error {
	return nil
}

func SetRedisSentinelSentinelsInfo(r map[string]string) error {
	return nil
}
