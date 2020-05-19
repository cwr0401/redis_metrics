package info

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type RedisStatsCollector struct {
	TotalConnectionsReceived   *prometheus.GaugeVec
	TotalCommandsProcessed     *prometheus.GaugeVec
	InstantaneousOpsPerSec     *prometheus.GaugeVec
	TotalNetInputBytes         *prometheus.GaugeVec
	TotalNetOutputBytes        *prometheus.GaugeVec
	InstantaneousInputKbps     *prometheus.GaugeVec
	InstantaneousOutputKbps    *prometheus.GaugeVec
	RejectedConnections        *prometheus.GaugeVec
	SyncFull                   *prometheus.GaugeVec
	SyncPartialOK              *prometheus.GaugeVec
	SyncPartialErr             *prometheus.GaugeVec
	ExpiredKeys                *prometheus.GaugeVec
	ExpiredStalePerc           *prometheus.GaugeVec
	ExpiredTimeCapReachedCount *prometheus.GaugeVec
	EvictedKeys                *prometheus.GaugeVec
	KeyspaceHits               *prometheus.GaugeVec
	KeyspaceMisses             *prometheus.GaugeVec
	PubsubChannels             *prometheus.GaugeVec
	PubsubPatterns             *prometheus.GaugeVec
	LatestForkUsec             *prometheus.GaugeVec
	MigrateCachedSockets       *prometheus.GaugeVec
	SlaveExpiresTrackedKeys    *prometheus.GaugeVec
	ActiveDefragHits           *prometheus.GaugeVec
	ActiveDefragMisses         *prometheus.GaugeVec
	ActiveDefragKeyHits        *prometheus.GaugeVec
	ActiveDefragKeyMisses      *prometheus.GaugeVec
}

func NewRedisStatsCollector() *RedisStatsCollector {

	var (
		//total_connections_received:2
		redisStatsTotalConnectionsReceived = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "total_connections_received",
			Help:      "Total number of connections accepted by the server.",
		}, []string{"node_name", "node_address"})

		//total_commands_processed:1
		redisStatsTotalCommandsProcessed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "total_commands_processed",
			Help:      "Total number of commands processed by the server.",
		}, []string{"node_name", "node_address"})

		//instantaneous_ops_per_sec:0
		redisStatsInstantaneousOpsPerSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "instantaneous_ops_per_sec",
			Help:      "Number of commands processed per second.",
		}, []string{"node_name", "node_address"})

		//total_net_input_bytes:90
		redisStatsTotalNetInputBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "total_net_input_bytes",
			Help:      "The total number of bytes read from the network.",
		}, []string{"node_name", "node_address"})

		//total_net_output_bytes:107
		redisStatsTotalNetOutputBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "total_net_output_bytes",
			Help:      "The total number of bytes written to the network.",
		}, []string{"node_name", "node_address"})

		//instantaneous_input_kbps:0.00
		redisStatsInstantaneousInputKbps = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "instantaneous_input_kbps",
			Help:      "The network's read rate per second in KB/sec.",
		}, []string{"node_name", "node_address"})

		//instantaneous_output_kbps:0.00
		redisStatsInstantaneousOutputKbps = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "instantaneous_output_kbps",
			Help:      "The network's write rate per second in KB/sec.",
		}, []string{"node_name", "node_address"})

		//rejected_connections:0
		redisStatsRejectedConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "rejected_connections",
			Help:      "Number of connections rejected because of maxclients limit.",
		}, []string{"node_name", "node_address"})

		//sync_full:0
		redisStatsSyncFull = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "sync_full",
			Help:      "The number of full resyncs with replicas.",
		}, []string{"node_name", "node_address"})

		//sync_partial_ok:0
		redisStatsSyncPartialOK = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "sync_partial_ok",
			Help:      "The number of accepted partial resync requests.",
		}, []string{"node_name", "node_address"})

		//sync_partial_err:0
		redisStatsSyncPartialErr = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "sync_partial_err",
			Help:      "The number of denied partial resync requests.",
		}, []string{"node_name", "node_address"})

		//expired_keys:0
		redisStatsExpiredKeys = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "expired_keys",
			Help:      "Total number of key expiration events.",
		}, []string{"node_name", "node_address"})

		//expired_stale_perc:0.00
		redisStatsExpiredStalePerc = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "expired_stale_perc",
			Help:      "",
		}, []string{"node_name", "node_address"})

		//expired_time_cap_reached_count:0
		redisStatsExpiredTimeCapReachedCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "expired_time_cap_reached_count",
			Help:      "",
		}, []string{"node_name", "node_address"})

		//evicted_keys:0
		redisStatsEvictedKeys = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "evicted_keys",
			Help:      "Number of evicted keys due to maxmemory limit.",
		}, []string{"node_name", "node_address"})

		//keyspace_hits:0
		redisStatsKeyspaceHits = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "keyspace_hits",
			Help:      "Number of successful lookup of keys in the main dictionary.",
		}, []string{"node_name", "node_address"})

		//keyspace_misses:0
		redisStatsKeyspaceMisses = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "keyspace_misses",
			Help:      "Number of failed lookup of keys in the main dictionary.",
		}, []string{"node_name", "node_address"})

		//pubsub_channels:0
		redisStatsPubsubChannels = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "pubsub_channels",
			Help:      "Global number of pub/sub channels with client subscriptions.",
		}, []string{"node_name", "node_address"})

		//pubsub_patterns:0
		redisStatsPubsubPatterns = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "pubsub_patterns",
			Help:      "Global number of pub/sub pattern with client subscriptions.",
		}, []string{"node_name", "node_address"})

		//latest_fork_usec:0
		redisStatsLatestForkUsec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "latest_fork_usec",
			Help:      "Duration of the latest fork operation in microseconds.",
		}, []string{"node_name", "node_address"})

		//migrate_cached_sockets:0
		redisStatsMigrateCachedSockets = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "migrate_cached_sockets",
			Help:      "The number of sockets open for MIGRATE purposes.",
		}, []string{"node_name", "node_address"})

		//slave_expires_tracked_keys:0
		redisStatsSlaveExpiresTrackedKeys = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "slave_expires_tracked_keys",
			Help:      "The number of keys tracked for expiry purposes (applicable only to writable replicas).",
		}, []string{"node_name", "node_address"})

		//active_defrag_hits:0
		redisStatsActiveDefragHits = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "active_defrag_hits",
			Help:      "Number of value reallocations performed by active the defragmentation process.",
		}, []string{"node_name", "node_address"})

		//active_defrag_misses:0
		redisStatsActiveDefragMisses = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "active_defrag_misses",
			Help:      "Number of aborted value reallocations started by the active defragmentation process.",
		}, []string{"node_name", "node_address"})

		//active_defrag_key_hits:0
		redisStatsActiveDefragKeyHits = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "active_defrag_key_hits",
			Help:      "Number of keys that were actively defragmented.",
		}, []string{"node_name", "node_address"})

		//active_defrag_key_misses:0
		redisStatsActiveDefragKeyMisses = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "stats",
			Name:      "active_defrag_key_misses",
			Help:      "Number of keys that were skipped by the active defragmentation process.",
		}, []string{"node_name", "node_address"})
	)
	return &RedisStatsCollector{
		redisStatsTotalConnectionsReceived,
		redisStatsTotalCommandsProcessed,
		redisStatsInstantaneousOpsPerSec,
		redisStatsTotalNetInputBytes,
		redisStatsTotalNetOutputBytes,
		redisStatsInstantaneousInputKbps,
		redisStatsInstantaneousOutputKbps,
		redisStatsRejectedConnections,
		redisStatsSyncFull,
		redisStatsSyncPartialOK,
		redisStatsSyncPartialErr,
		redisStatsExpiredKeys,
		redisStatsExpiredStalePerc,
		redisStatsExpiredTimeCapReachedCount,
		redisStatsEvictedKeys,
		redisStatsKeyspaceHits,
		redisStatsKeyspaceMisses,
		redisStatsPubsubChannels,
		redisStatsPubsubPatterns,
		redisStatsLatestForkUsec,
		redisStatsMigrateCachedSockets,
		redisStatsSlaveExpiresTrackedKeys,
		redisStatsActiveDefragHits,
		redisStatsActiveDefragMisses,
		redisStatsActiveDefragKeyHits,
		redisStatsActiveDefragKeyMisses,
	}
}

func (m *RedisStatsCollector) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(m.TotalConnectionsReceived)
	registry.MustRegister(m.TotalCommandsProcessed)
	registry.MustRegister(m.InstantaneousOpsPerSec)
	registry.MustRegister(m.TotalNetInputBytes)
	registry.MustRegister(m.TotalNetOutputBytes)
	registry.MustRegister(m.InstantaneousInputKbps)
	registry.MustRegister(m.InstantaneousOutputKbps)
	registry.MustRegister(m.RejectedConnections)
	registry.MustRegister(m.SyncFull)
	registry.MustRegister(m.SyncPartialOK)
	registry.MustRegister(m.SyncPartialErr)
	registry.MustRegister(m.ExpiredKeys)
	registry.MustRegister(m.ExpiredStalePerc)
	registry.MustRegister(m.ExpiredTimeCapReachedCount)
	registry.MustRegister(m.EvictedKeys)
	registry.MustRegister(m.KeyspaceHits)
	registry.MustRegister(m.KeyspaceMisses)
	registry.MustRegister(m.PubsubChannels)
	registry.MustRegister(m.PubsubPatterns)
	registry.MustRegister(m.LatestForkUsec)
	registry.MustRegister(m.MigrateCachedSockets)
	registry.MustRegister(m.SlaveExpiresTrackedKeys)
	registry.MustRegister(m.ActiveDefragHits)
	registry.MustRegister(m.ActiveDefragMisses)
	registry.MustRegister(m.ActiveDefragKeyHits)
	registry.MustRegister(m.ActiveDefragKeyMisses)

}

func (m *RedisStatsCollector) Unregister(registry *prometheus.Registry) bool {
	if !registry.Unregister(m.TotalConnectionsReceived) {
		return false
	}
	if !registry.Unregister(m.TotalCommandsProcessed) {
		return false
	}
	if !registry.Unregister(m.InstantaneousOpsPerSec) {
		return false
	}
	if !registry.Unregister(m.TotalNetInputBytes) {
		return false
	}
	if !registry.Unregister(m.TotalNetOutputBytes) {
		return false
	}
	if !registry.Unregister(m.InstantaneousInputKbps) {
		return false
	}
	if !registry.Unregister(m.InstantaneousOutputKbps) {
		return false
	}
	if !registry.Unregister(m.RejectedConnections) {
		return false
	}
	if !registry.Unregister(m.SyncFull) {
		return false
	}
	if !registry.Unregister(m.SyncPartialOK) {
		return false
	}
	if !registry.Unregister(m.SyncPartialErr) {
		return false
	}
	if !registry.Unregister(m.ExpiredKeys) {
		return false
	}
	if !registry.Unregister(m.ExpiredStalePerc) {
		return false
	}
	if !registry.Unregister(m.ExpiredTimeCapReachedCount) {
		return false
	}
	if !registry.Unregister(m.EvictedKeys) {
		return false
	}
	if !registry.Unregister(m.KeyspaceHits) {
		return false
	}
	if !registry.Unregister(m.KeyspaceMisses) {
		return false
	}
	if !registry.Unregister(m.PubsubChannels) {
		return false
	}
	if !registry.Unregister(m.PubsubPatterns) {
		return false
	}
	if !registry.Unregister(m.LatestForkUsec) {
		return false
	}
	if !registry.Unregister(m.MigrateCachedSockets) {
		return false
	}
	if !registry.Unregister(m.SlaveExpiresTrackedKeys) {
		return false
	}
	if !registry.Unregister(m.ActiveDefragHits) {
		return false
	}
	if !registry.Unregister(m.ActiveDefragMisses) {
		return false
	}
	if !registry.Unregister(m.ActiveDefragKeyHits) {
		return false
	}
	if !registry.Unregister(m.ActiveDefragKeyMisses) {
		return false
	}
	return true
}

func (m *RedisStatsCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	// total_connections_received:33
	if totalConnectionsReceivedStr, ok := r["total_connections_received"]; ok {
		if totalConnectionsReceived, err := strconv.Atoi(totalConnectionsReceivedStr); err == nil {
			m.TotalConnectionsReceived.WithLabelValues(nodeName, nodeAddress).Set(
				float64(totalConnectionsReceived))
		}
	}
	//total_commands_processed:96955
	if totalCommandsProcessedStr, ok := r["total_commands_processed"]; ok {
		if totalCommandsProcessed, err := strconv.Atoi(totalCommandsProcessedStr); err == nil {
			m.TotalCommandsProcessed.WithLabelValues(nodeName, nodeAddress).Set(float64(totalCommandsProcessed))
		}
	}
	//instantaneous_ops_per_sec:7
	if instantaneousOpsPerSecStr, ok := r["instantaneous_ops_per_sec"]; ok {
		if instantaneousOpsPerSec, err := strconv.Atoi(instantaneousOpsPerSecStr); err == nil {
			m.InstantaneousOpsPerSec.WithLabelValues(nodeName, nodeAddress).Set(float64(instantaneousOpsPerSec))
		}
	}
	//total_net_input_bytes:4695881
	if totalNetInputBytesStr, ok := r["total_net_input_bytes"]; ok {
		if totalNetInputBytes, err := strconv.Atoi(totalNetInputBytesStr); err == nil {
			m.TotalNetInputBytes.WithLabelValues(nodeName, nodeAddress).Set(float64(totalNetInputBytes))
		}
	}
	//total_net_output_bytes:30079574
	if totalNetOutputBytesStr, ok := r["total_net_output_bytes"]; ok {
		if totalNetOutputBytes, err := strconv.Atoi(totalNetOutputBytesStr); err == nil {
			m.TotalNetOutputBytes.WithLabelValues(nodeName, nodeAddress).Set(float64(totalNetOutputBytes))
		}
	}
	//instantaneous_input_kbps:0.31
	if instantaneousInputKbpsStr, ok := r["instantaneous_input_kbps"]; ok {
		if instantaneousInputKbps, err := strconv.ParseFloat(instantaneousInputKbpsStr, 64); err == nil {
			m.InstantaneousInputKbps.WithLabelValues(nodeName, nodeAddress).Set(instantaneousInputKbps)
		}
	}
	//instantaneous_output_kbps:7.70
	if instantaneusOutputKbpsStr, ok := r["instantaneous_output_kbps"]; ok {
		if instantaneusOutputKbps, err := strconv.ParseFloat(instantaneusOutputKbpsStr, 64); err == nil {
			m.InstantaneousOutputKbps.WithLabelValues(nodeName, nodeAddress).Set(instantaneusOutputKbps)
		}
	}
	//rejected_connections:0
	if rejectedConnectionsStr, ok := r["rejected_connections"]; ok {
		if rejectedConnections, err := strconv.Atoi(rejectedConnectionsStr); err == nil {
			m.RejectedConnections.WithLabelValues(nodeName, nodeAddress).Set(float64(rejectedConnections))
		}
	}
	//sync_full:2
	if syncFullStr, ok := r["sync_full"]; ok {
		if syncFull, err := strconv.Atoi(syncFullStr); err == nil {
			m.SyncFull.WithLabelValues(nodeName, nodeAddress).Set(float64(syncFull))
		}
	}
	//sync_partial_ok:2
	if syncPartialOKStr, ok := r["sync_partial_ok"]; ok {
		if syncPartialOK, err := strconv.Atoi(syncPartialOKStr); err == nil {
			m.SyncPartialOK.WithLabelValues(nodeName, nodeAddress).Set(float64(syncPartialOK))
		}
	}
	//sync_partial_err:2
	if syncPartialERRStr, ok := r["sync_partial_err"]; ok {
		if syncPartialERR, err := strconv.Atoi(syncPartialERRStr); err == nil {
			m.SyncPartialErr.WithLabelValues(nodeName, nodeAddress).Set(float64(syncPartialERR))
		}
	}
	//expired_keys:0
	if expiredKeysStr, ok := r["expired_keys"]; ok {
		if expiredKeys, err := strconv.Atoi(expiredKeysStr); err == nil {
			m.ExpiredKeys.WithLabelValues(nodeName, nodeAddress).Set(float64(expiredKeys))
		}
	}
	//expired_stale_perc:0.00
	if expiredStalePercStr, ok := r["expired_stale_perc"]; ok {
		if expiredStalePerc, err := strconv.ParseFloat(expiredStalePercStr, 64); err == nil {
			m.ExpiredStalePerc.WithLabelValues(nodeName, nodeAddress).Set(expiredStalePerc)
		}
	}
	//expired_time_cap_reached_count:0
	if expiredTimeCapReadchedCountStr, ok := r["expired_time_cap_reached_count"]; ok {
		if expiredTimeCapReadchedCount, err := strconv.Atoi(expiredTimeCapReadchedCountStr); err == nil {
			m.ExpiredTimeCapReachedCount.WithLabelValues(nodeName, nodeAddress).Set(
				float64(expiredTimeCapReadchedCount))
		}
	}
	//evicted_keys:0
	if evictedKeysStr, ok := r["evicted_keys"]; ok {
		if evictedKeys, err := strconv.Atoi(evictedKeysStr); err == nil {
			m.EvictedKeys.WithLabelValues(nodeName, nodeAddress).Set(float64(evictedKeys))
		}
	}
	//keyspace_hits:0
	if keyspaceHitsStr, ok := r["keyspace_hits"]; ok {
		if keyspaceHits, err := strconv.Atoi(keyspaceHitsStr); err == nil {
			m.KeyspaceHits.WithLabelValues(nodeName, nodeAddress).Set(float64(keyspaceHits))
		}
	}
	//keyspace_misses:0
	if keyspaceMissesStr, ok := r["keyspace_misses"]; ok {
		if keyspaceMisses, err := strconv.Atoi(keyspaceMissesStr); err == nil {
			m.KeyspaceMisses.WithLabelValues(nodeName, nodeAddress).Set(float64(keyspaceMisses))
		}
	}
	//pubsub_channels:1
	if pubsubChannelsStr, ok := r["pubsub_channels"]; ok {
		if pubsubChannels, err := strconv.Atoi(pubsubChannelsStr); err == nil {
			m.PubsubChannels.WithLabelValues(nodeName, nodeAddress).Set(float64(pubsubChannels))
		}
	}
	//pubsub_patterns:0
	if pubsubPatternsStr, ok := r["pubsub_patterns"]; ok {
		if pubsubPatterns, err := strconv.Atoi(pubsubPatternsStr); err == nil {
			m.PubsubChannels.WithLabelValues(nodeName, nodeAddress).Set(float64(pubsubPatterns))
		}
	}
	//latest_fork_usec:366
	if latestForkUsecStr, ok := r["latest_fork_usec"]; ok {
		if latestForkUsec, err := strconv.Atoi(latestForkUsecStr); err == nil {
			m.LatestForkUsec.WithLabelValues(nodeName, nodeAddress).Set(float64(latestForkUsec))
		}
	}
	//migrate_cached_sockets:0
	if migrateCachedSocketsStr, ok := r["migrate_cached_sockets"]; ok {
		if migrateCachedSockets, err := strconv.Atoi(migrateCachedSocketsStr); err == nil {
			m.MigrateCachedSockets.WithLabelValues(nodeName, nodeAddress).Set(float64(migrateCachedSockets))
		}
	}
	//slave_expires_tracked_keys:0
	if slaveExpiresTrackedKeysStr, ok := r["slave_expires_tracked_keys"]; ok {
		if slaveExpiresTrackedKeys, err := strconv.Atoi(slaveExpiresTrackedKeysStr); err == nil {
			m.SlaveExpiresTrackedKeys.WithLabelValues(nodeName, nodeAddress).Set(
				float64(slaveExpiresTrackedKeys))
		}
	}
	//active_defrag_hits:0
	if activeDefragHitsStr, ok := r["active_defrag_hits"]; ok {
		if activeDefragHits, err := strconv.Atoi(activeDefragHitsStr); err == nil {
			m.ActiveDefragHits.WithLabelValues(nodeName, nodeAddress).Set(float64(activeDefragHits))
		}
	}
	//active_defrag_misses:0
	if activeDefragMissesStr, ok := r["active_defrag_misses"]; ok {
		if activeDefragMisses, err := strconv.Atoi(activeDefragMissesStr); err == nil {
			m.ActiveDefragMisses.WithLabelValues(nodeName, nodeAddress).Set(float64(activeDefragMisses))
		}
	}
	//active_defrag_key_hits:0
	if activeDefragKeyHitsStr, ok := r["active_defrag_key_hits"]; ok {
		if activeDefragKeyHits, err := strconv.Atoi(activeDefragKeyHitsStr); err == nil {
			m.ActiveDefragKeyHits.WithLabelValues(nodeName, nodeAddress).Set(float64(activeDefragKeyHits))
		}
	}
	//active_defrag_key_misses:0
	if activeDefragKeyMissesStr, ok := r["active_defrag_key_misses"]; ok {
		if activeDefragKeyMisses, err := strconv.Atoi(activeDefragKeyMissesStr); err == nil {
			m.ActiveDefragKeyMisses.WithLabelValues(nodeName, nodeAddress).Set(float64(activeDefragKeyMisses))
		}
	}
	return nil
}
