package info

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Server
	// up
	RedisServerUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "server",
		Name:      "up",
		Help:      "Value is 1 if Redis server alive, 0 otherwise.",
	},
		[]string{"node_name", "node_address"})

	// info
	//redis_version:5.0.7
	//redis_git_sha1:00000000
	//redis_git_dirty:0
	//redis_build_id:7359662505fc6f11
	//redis_mode:standalone
	//os:Linux 4.9.184-linuxkit x86_64
	//arch_bits:64
	//multiplexing_api:epoll
	//atomicvar_api:atomic-builtin
	//gcc_version:8.3.0
	//process_id:1
	//run_id:567eb7e45beb623b95f0785d8841a1824571a232
	//tcp_port:6379
	//uptime_in_seconds:358
	//uptime_in_days:0
	//hz:10
	//configured_hz:10
	//lru_clock:14370486
	//executable:/data/redis-server
	//config_file:/etc/redis.conf

	// redis_version: Version of the Redis server
	// redis_git_sha1: Git SHA1
	// redis_git_dirty: Git dirty flag
	// redis_build_id: The build id
	// redis_mode: The server's mode ("standalone", "sentinel" or "cluster")
	// os: Operating system hosting the Redis server
	// arch_bits: Architecture (32 or 64 bits)
	// multiplexing_api: Event loop mechanism used by Redis
	// atomicvar_api: Atomicvar API used by Redis
	// gcc_version: Version of the GCC compiler used to compile the Redis server
	// process_id: PID of the server process
	// run_id: Random value identifying the Redis server (to be used by Sentinel and Cluster)
	// tcp_port: TCP/IP listen port
	// uptime_in_seconds: Number of seconds since Redis server start
	// uptime_in_days: Same value expressed in days
	// hz: The server's frequency setting
	// lru_clock: Clock incrementing every minute, for LRU management
	// executable: The path to the server's executable
	// config_file: The path to the config file

	RedisServerInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "server",
			Name:      "info",
			Help:      "Information about the Redis server.",
		},
		[]string{
			"node_name",
			"node_address",
			"redis_version",
			"redis_git_sha1",
			"redis_git_dirty",
			"redis_build_id",
			"redis_mode",
			"os",
			"arch_bits",
			"multiplexing_api",
			"atomicvar_api",
			"gcc_version",
			// "process_id",
			// "run_id",
			"tcp_port",
			"executable",
			"config_file",
		})

	// uptime_in_seconds
	RedisServerUptimeInSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "server",
		Name:      "uptime_in_seconds",
		Help:      "Number of seconds since Redis server start.",
	},
		[]string{"node_name", "node_address"})

	// uptime_in_days
	RedisServerUptimeInDays = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "server",
		Name:      "uptime_in_days",
		Help:      "Number of days since Redis server start.",
	},
		[]string{"node_name", "node_address"})

	// hz
	RedisServerHz = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "server",
		Name:      "hz",
		Help:      "The server's frequency setting.",
	},
		[]string{"node_name", "node_address"})

	// configured_hz
	RedisServerConfiguredHz = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "server",
		Name:      "configured_hz",
		Help:      "The server's configured frequency setting.",
	},
		[]string{"node_name", "node_address"})

	// lru_clock
	RedisServerLruClock = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "server",
		Name:      "lru_clock",
		Help:      "Clock incrementing every minute, for LRU management.",
	},
		[]string{"node_name", "node_address"})
)

func SetRedisServer(nodeName, nodeAddress string, r map[string]string) error {
	redisVersion, ok := r["redis_version"]
	if !ok {
		redisVersion = "unknown"
	}
	redisGitSha1, ok := r["redis_git_sha1"]
	if !ok {
		redisGitSha1 = "unknown"
	}
	redisGitDirty, ok := r["redis_git_dirty"]
	if !ok {
		redisGitDirty = "unknown"
	}
	redisBuildId, ok := r["redis_build_id"]
	if !ok {
		redisBuildId = "unknown"
	}
	redisMode, ok := r["redis_mode"]
	if !ok {
		redisMode = "standalone"
	}
	os, ok := r["os"]
	if !ok {
		os = "unknown"
	}
	archBits, ok := r["arch_bits"]
	if !ok {
		archBits = "unknown"
	}
	multiplexingAPI, ok := r["multiplexing_api"]
	if !ok {
		multiplexingAPI = "unknown"
	}
	atomicvarAPI, ok := r["atomicvar_api"]
	if !ok {
		atomicvarAPI = "unknown"
	}
	gccVersion, ok := r["gcc_version"]
	if !ok {
		gccVersion = "unknown"
	}
	//processId, ok := r["process_id"]
	//if !ok {
	//	processId = "unknown"
	//}
	//runId, ok := r["run_id"]
	//if !ok {
	//	runId = "unknown"
	//}
	tcpPort, ok := r["tcp_port"]
	if !ok {
		tcpPort = "unknown"
	}
	executable, ok := r["executable"]
	if !ok {
		executable = "unknown"
	}
	configFile, ok := r["config_file"]
	if !ok {
		configFile = "unknown"
	}

	RedisServerInfo.WithLabelValues(
		nodeName, nodeAddress, redisVersion, redisGitSha1, redisGitDirty, redisBuildId, redisMode, os,
		archBits, multiplexingAPI, atomicvarAPI, gccVersion, tcpPort, executable, configFile,
	).Set(1)

	if uptimeInSecondsStr, ok := r["uptime_in_seconds"]; ok {
		if uptimeInSeconds, err := strconv.Atoi(uptimeInSecondsStr); err == nil {
			RedisServerUptimeInSeconds.WithLabelValues(nodeName, nodeAddress).Set(float64(uptimeInSeconds))
		}
	}

	if uptimeInDaysStr, ok := r["uptime_in_days"]; ok {
		if uptimeInDays, err := strconv.Atoi(uptimeInDaysStr); err == nil {
			RedisServerUptimeInDays.WithLabelValues(nodeName, nodeAddress).Set(float64(uptimeInDays))
		}
	}

	if hzStr, ok := r["hz"]; ok {
		if hz, err := strconv.Atoi(hzStr); err == nil {
			RedisServerHz.WithLabelValues(nodeName, nodeAddress).Set(float64(hz))
		}
	}

	if configuredHzStr, ok := r["configured_hz"]; ok {
		if configuredHz, err := strconv.Atoi(configuredHzStr); err == nil {
			RedisServerConfiguredHz.WithLabelValues(nodeName, nodeAddress).Set(float64(configuredHz))
		}
	}

	if lruClockStr, ok := r["lru_clock"]; ok {
		if lruClock, err := strconv.Atoi(lruClockStr); err == nil {
			RedisServerLruClock.WithLabelValues(nodeName, nodeAddress).Set(float64(lruClock))
		}
	}
	return nil
}
