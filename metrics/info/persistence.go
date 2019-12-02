package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var (
	// Persistence
	//loading:0
	//rdb_changes_since_last_save:0
	//rdb_bgsave_in_progress:0
	//rdb_last_save_time:1574651217
	//rdb_last_bgsave_status:ok
	//rdb_last_bgsave_time_sec:-1
	//rdb_current_bgsave_time_sec:-1
	//rdb_last_cow_size:0
	//aof_enabled:0
	//aof_rewrite_in_progress:0
	//aof_rewrite_scheduled:0
	//aof_last_rewrite_time_sec:-1
	//aof_current_rewrite_time_sec:-1
	//aof_last_bgrewrite_status:ok
	//aof_last_write_status:ok
	//aof_last_cow_size:0

	// loading
	RedisPersistenceLoading = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "loading",
		Help:      "Flag indicating if the load of a dump file is on-going.",
	}, []string{"node_name", "node_address"})

	// rdb_changes_since_last_save
	RedisPersistenceRdbChangesSinceLastSave = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "rdb_changes_since_last_save",
		Help:      "Number of changes since the last dump.",
	}, []string{"node_name", "node_address"})

	// rdb_bgsave_in_progress
	RedisPersistenceRdbBgsaveInProgress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "rdb_bgsave_in_progress",
		Help:      "Flag indicating a RDB save is on-going.",
	}, []string{"node_name", "node_address"})

	// rdb_last_save_time
	RedisPersistenceRdbLastSaveTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "rdb_last_save_time",
		Help:      "Epoch-based timestamp of last successful RDB save.",
	}, []string{"node_name", "node_address"})

	// rdb_last_bgsave_status
	RedisPersistenceRdbLastBgsaveStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "rdb_last_bgsave_status",
		Help:      "Status of the last RDB save operation.",
	}, []string{"node_name", "node_address"})

	// rdb_last_bgsave_time_sec
	RedisPersistenceRdbLastBgsaveTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "rdb_last_bgsave_time_sec",
		Help:      "Duration of the last RDB save operation in seconds",
	}, []string{"node_name", "node_address"})

	// rdb_current_bgsave_time_sec
	RedisPersistenceRdbCurrentBgsaveTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "rdb_current_bgsave_time_sec",
		Help:      "Duration of the on-going RDB save operation if any",
	}, []string{"node_name", "node_address"})

	// rdb_last_cow_size
	RedisPersistenceRdbLastCowSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "rdb_last_cow_size",
		Help:      "The size in bytes of copy-on-write allocations during the last RDB save operation",
	}, []string{"node_name", "node_address"})

	// aof_enabled
	RedisPersistenceAofEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_enabled",
		Help:      "Flag indicating AOF logging is activated",
	}, []string{"node_name", "node_address"})

	// aof_rewrite_in_progress
	RedisPersistenceAofRewriteInProgress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_rewrite_in_progress",
		Help:      "Flag indicating a AOF rewrite operation is on-going",
	}, []string{"node_name", "node_address"})

	// aof_rewrite_scheduled
	RedisPersistenceAofRewriteScheduled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_rewrite_scheduled",
		Help:      "Flag indicating an AOF rewrite operation will be scheduled once the on-going RDB save is complete.",
	}, []string{"node_name", "node_address"})

	// aof_last_rewrite_time_sec
	RedisPersistenceAofLastRewriteTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_last_rewrite_time_sec",
		Help:      "Duration of the last AOF rewrite operation in seconds",
	}, []string{"node_name", "node_address"})

	// aof_current_rewrite_time_sec
	RedisPersistenceAofCurrentRewriteTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_current_rewrite_time_sec",
		Help:      "Duration of the on-going AOF rewrite operation if any",
	}, []string{"node_name", "node_address"})

	// aof_last_bgrewrite_status
	RedisPersistenceAofLastBgrewriteStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_last_bgrewrite_status",
		Help:      "Status of the last AOF rewrite operation",
	}, []string{"node_name", "node_address"})

	// aof_last_write_status
	RedisPersistenceAofLastWriteStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_last_write_status",
		Help:      "Status of the last write operation to the AOF",
	}, []string{"node_name", "node_address"})

	// aof_last_cow_size
	RedisPersistenceAofLastCowSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "persistence",
		Name:      "aof_last_cow_size",
		Help:      "The size in bytes of copy-on-write allocations during the last AOF rewrite operation",
	}, []string{"node_name", "node_address"})
)

func SetRedisPersistence(nodeName, nodeAddress string, r map[string]string) error {
	if loadingStr, ok := r["loading"]; ok {
		if loading, err := strconv.Atoi(loadingStr); err == nil {
			RedisPersistenceLoading.WithLabelValues(nodeName, nodeAddress).Set(float64(loading))
		}
	}
	if rdbChangersSinceLastSaveStr, ok := r["rdb_changes_since_last_save"]; ok {
		if rdbChangersSinceLastSave, err := strconv.Atoi(rdbChangersSinceLastSaveStr); err == nil {
			RedisPersistenceRdbChangesSinceLastSave.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbChangersSinceLastSave))
		}
	}
	if rdbBgsaveInProgressStr, ok := r["rdb_bgsave_in_progress"]; ok {
		if rdbBgsaveInProgress, err := strconv.Atoi(rdbBgsaveInProgressStr); err == nil {
			RedisPersistenceRdbBgsaveInProgress.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbBgsaveInProgress))
		}
	}
	if rdbLastSaveTimeStr, ok := r["rdb_last_save_time"]; ok {
		if rdbLastSaveTime, err := strconv.Atoi(rdbLastSaveTimeStr); err == nil {
			RedisPersistenceRdbLastSaveTime.WithLabelValues(nodeName, nodeAddress).Set(float64(rdbLastSaveTime))
		}
	}
	if rdbLastBgsaveStatus, ok := r["rdb_last_bgsave_status"]; ok {
		switch rdbLastBgsaveStatus {
		case "ok":
			RedisPersistenceRdbLastBgsaveStatus.WithLabelValues(nodeName, nodeAddress).Set(1)
		default:
			RedisPersistenceRdbLastBgsaveStatus.WithLabelValues(nodeName, nodeAddress).Set(0)
		}
	}
	if rdbLastBgsaveTimeSecStr, ok := r["rdb_last_bgsave_time_sec"]; ok {
		if rdbLastBgsaveTimeSec, err := strconv.Atoi(rdbLastBgsaveTimeSecStr); err == nil {
			RedisPersistenceRdbLastBgsaveTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbLastBgsaveTimeSec))
		}
	}
	if rdbCurrentBgsaveTimeSecStr, ok := r["rdb_current_bgsave_time_sec"]; ok {
		if rdbCurrentBgsaveTimeSec, err := strconv.Atoi(rdbCurrentBgsaveTimeSecStr); err == nil {
			RedisPersistenceRdbCurrentBgsaveTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbCurrentBgsaveTimeSec))
		}
	}
	if rdbLastCowSizeStr, ok := r["rdb_last_cow_size"]; ok {
		if rdbLastCowSize, err := strconv.Atoi(rdbLastCowSizeStr); err == nil {
			RedisPersistenceRdbLastCowSize.WithLabelValues(nodeName, nodeAddress).Set(float64(rdbLastCowSize))
		}
	}
	if aofEnabledStr, ok := r["aof_enabled"]; ok {
		if aofEnabled, err := strconv.Atoi(aofEnabledStr); err == nil {
			RedisPersistenceAofEnabled.WithLabelValues(nodeName, nodeAddress).Set(float64(aofEnabled))
		}
	}
	if aofRewriteInProgressStr, ok := r["aof_rewrite_in_progress"]; ok {
		if aofRewriteInProgress, err := strconv.Atoi(aofRewriteInProgressStr); err == nil {
			RedisPersistenceAofRewriteInProgress.WithLabelValues(nodeName, nodeAddress).Set(
				float64(aofRewriteInProgress))
		}
	}
	if aofRewriteScheduledStr, ok := r["aof_rewrite_scheduled"]; ok {
		if aofRewriteScheduled, err := strconv.Atoi(aofRewriteScheduledStr); err == nil {
			RedisPersistenceAofRewriteScheduled.WithLabelValues(nodeName, nodeAddress).Set(float64(aofRewriteScheduled))
		}
	}
	if aofLastRewriteTimeSecStr, ok := r["aof_last_rewrite_time_sec"]; ok {
		if aofLastRewriteTimeSec, err := strconv.Atoi(aofLastRewriteTimeSecStr); err == nil {
			RedisPersistenceAofLastRewriteTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(aofLastRewriteTimeSec))
		}
	}
	if aofCurrentRewriteTimeSecStr, ok := r["aof_current_rewrite_time_sec"]; ok {
		if aofCurrentRewriteTimeSec, err := strconv.Atoi(aofCurrentRewriteTimeSecStr); err == nil {
			RedisPersistenceAofCurrentRewriteTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(aofCurrentRewriteTimeSec))
		}
	}
	if aofLastBgrewriteStatus, ok := r["aof_last_bgrewrite_status"]; ok {
		switch aofLastBgrewriteStatus {
		case "ok":
			RedisPersistenceAofLastBgrewriteStatus.WithLabelValues(nodeName, nodeAddress).Set(1)
		default:
			RedisPersistenceAofLastBgrewriteStatus.WithLabelValues(nodeName, nodeAddress).Set(0)
		}
	}
	if aofLastWriteStatus, ok := r["aof_last_write_status"]; ok {
		switch aofLastWriteStatus {
		case "ok":
			RedisPersistenceAofLastWriteStatus.WithLabelValues(nodeName, nodeAddress).Set(1)
		default:
			RedisPersistenceAofLastWriteStatus.WithLabelValues(nodeName, nodeAddress).Set(0)
		}
	}
	if aofLastCowSizeStr, ok := r["aof_last_cow_size"]; ok {
		if aofLastCowSize, err := strconv.Atoi(aofLastCowSizeStr); err == nil {
			RedisPersistenceAofLastCowSize.WithLabelValues(nodeName, nodeAddress).Set(float64(aofLastCowSize))
		}
	}
	return nil
}
