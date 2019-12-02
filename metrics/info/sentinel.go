package info

import "github.com/prometheus/client_golang/prometheus"

var (
	// Sentinel
	RedisSentinelUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "up",
		Help:      "Value is 1 if Redis server alive, 0 otherwise.",
	},
		[]string{"node_name", "node_address"})

	// sentinel_masters:1
	RedisSentinelSentinelMasters = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "sentinel_masters",
		Help:      "",
	},
		[]string{"node_name", "node_address"})

	//sentinel_tilt:0
	RedisSentinelSentinelTilt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "sentinel_tilt",
		Help:      "",
	},
		[]string{"node_name", "node_address"})

	//sentinel_running_scripts:0
	RedisSentinelSentinelRunningScripts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "sentinel_running_scripts",
		Help:      "",
	},
		[]string{"node_name", "node_address"})

	//sentinel_scripts_queue_length:0
	RedisSentinelSentinelScriptsQueueLength = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "sentinel_scripts_queue_length",
		Help:      "",
	},
		[]string{"node_name", "node_address"})

	//sentinel_simulate_failure_flags:0
	RedisSentinelSentinelSimulateFailureFlags = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "sentinel_simulate_failure_flags",
		Help:      "",
	},
		[]string{"node_name", "node_address"})

	//master0:name=mymaster,status=ok,address=172.25.1.3:6379,slaves=2,sentinels=3
	RedisSentinelMaster = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "master",
		Help:      "",
	},
		[]string{"node_name", "node_address", "index", "name", "status", "address", "slaves", "sentinels"})

	// sentinel master
	RedisSentinelMastersInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
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

	RedisSentinelSlavesInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "slaves_info",
		Help:      "",
	},
		[]string{})

	RedisSentinelSentinelsInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "sentinel",
		Name:      "sentinels_info",
		Help:      "",
	},
		[]string{})
)

func SetRedisSentinel(nodeName, nodeAddress string, r map[string]string) error {
	return nil
}

func SetRedisSentinelMastersInfo(r map[string]string) error {
	RedisSentinelMastersInfo.WithLabelValues(r["name"], r["ip"], r["port"], r["runid"], r["flags"]).Inc()

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
