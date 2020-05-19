package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type RedisPersistenceCollector struct {
	Loading                  *prometheus.GaugeVec
	RdbChangesSinceLastSave  *prometheus.GaugeVec
	RdbBgsaveInProgress      *prometheus.GaugeVec
	RdbLastSaveTime          *prometheus.GaugeVec
	RdbLastBgsaveStatus      *prometheus.GaugeVec
	RdbLastBgsaveTimeSec     *prometheus.GaugeVec
	RdbCurrentBgsaveTimeSec  *prometheus.GaugeVec
	RdbLastCowSize           *prometheus.GaugeVec
	AofEnabled               *prometheus.GaugeVec
	AofRewriteInProgress     *prometheus.GaugeVec
	AofRewriteScheduled      *prometheus.GaugeVec
	AofLastRewriteTimeSec    *prometheus.GaugeVec
	AofCurrentRewriteTimeSec *prometheus.GaugeVec
	AofLastBgrewriteStatus   *prometheus.GaugeVec
	AofLastWriteStatus       *prometheus.GaugeVec
	AofLastCowSize           *prometheus.GaugeVec
}

func NewRedisPersistenceCollector() *RedisPersistenceCollector {
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
		redisPersistenceLoading = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "loading",
			Help:      "Flag indicating if the load of a dump file is on-going.",
		}, []string{"node_name", "node_address"})

		// rdb_changes_since_last_save
		redisPersistenceRdbChangesSinceLastSave = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "rdb_changes_since_last_save",
			Help:      "Number of changes since the last dump.",
		}, []string{"node_name", "node_address"})

		// rdb_bgsave_in_progress
		redisPersistenceRdbBgsaveInProgress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "rdb_bgsave_in_progress",
			Help:      "Flag indicating a RDB save is on-going.",
		}, []string{"node_name", "node_address"})

		// rdb_last_save_time
		redisPersistenceRdbLastSaveTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "rdb_last_save_time",
			Help:      "Epoch-based timestamp of last successful RDB save.",
		}, []string{"node_name", "node_address"})

		// rdb_last_bgsave_status
		redisPersistenceRdbLastBgsaveStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "rdb_last_bgsave_status",
			Help:      "Status of the last RDB save operation.",
		}, []string{"node_name", "node_address"})

		// rdb_last_bgsave_time_sec
		redisPersistenceRdbLastBgsaveTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "rdb_last_bgsave_time_sec",
			Help:      "Duration of the last RDB save operation in seconds",
		}, []string{"node_name", "node_address"})

		// rdb_current_bgsave_time_sec
		redisPersistenceRdbCurrentBgsaveTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "rdb_current_bgsave_time_sec",
			Help:      "Duration of the on-going RDB save operation if any",
		}, []string{"node_name", "node_address"})

		// rdb_last_cow_size
		redisPersistenceRdbLastCowSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "rdb_last_cow_size",
			Help:      "The size in bytes of copy-on-write allocations during the last RDB save operation",
		}, []string{"node_name", "node_address"})

		// aof_enabled
		redisPersistenceAofEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_enabled",
			Help:      "Flag indicating AOF logging is activated",
		}, []string{"node_name", "node_address"})

		// aof_rewrite_in_progress
		redisPersistenceAofRewriteInProgress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_rewrite_in_progress",
			Help:      "Flag indicating a AOF rewrite operation is on-going",
		}, []string{"node_name", "node_address"})

		// aof_rewrite_scheduled
		redisPersistenceAofRewriteScheduled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_rewrite_scheduled",
			Help:      "Flag indicating an AOF rewrite operation will be scheduled once the on-going RDB save is complete.",
		}, []string{"node_name", "node_address"})

		// aof_last_rewrite_time_sec
		redisPersistenceAofLastRewriteTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_last_rewrite_time_sec",
			Help:      "Duration of the last AOF rewrite operation in seconds",
		}, []string{"node_name", "node_address"})

		// aof_current_rewrite_time_sec
		redisPersistenceAofCurrentRewriteTimeSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_current_rewrite_time_sec",
			Help:      "Duration of the on-going AOF rewrite operation if any",
		}, []string{"node_name", "node_address"})

		// aof_last_bgrewrite_status
		redisPersistenceAofLastBgrewriteStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_last_bgrewrite_status",
			Help:      "Status of the last AOF rewrite operation",
		}, []string{"node_name", "node_address"})

		// aof_last_write_status
		redisPersistenceAofLastWriteStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_last_write_status",
			Help:      "Status of the last write operation to the AOF",
		}, []string{"node_name", "node_address"})

		// aof_last_cow_size
		redisPersistenceAofLastCowSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "persistence",
			Name:      "aof_last_cow_size",
			Help:      "The size in bytes of copy-on-write allocations during the last AOF rewrite operation",
		}, []string{"node_name", "node_address"})
	)
	return &RedisPersistenceCollector{
		redisPersistenceLoading,
		redisPersistenceRdbChangesSinceLastSave,
		redisPersistenceRdbBgsaveInProgress,
		redisPersistenceRdbLastSaveTime,
		redisPersistenceRdbLastBgsaveStatus,
		redisPersistenceRdbLastBgsaveTimeSec,
		redisPersistenceRdbCurrentBgsaveTimeSec,
		redisPersistenceRdbLastCowSize,
		redisPersistenceAofEnabled,
		redisPersistenceAofRewriteInProgress,
		redisPersistenceAofRewriteScheduled,
		redisPersistenceAofLastRewriteTimeSec,
		redisPersistenceAofCurrentRewriteTimeSec,
		redisPersistenceAofLastBgrewriteStatus,
		redisPersistenceAofLastWriteStatus,
		redisPersistenceAofLastCowSize,
	}
}

func (m *RedisPersistenceCollector) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(m.Loading)
	registry.MustRegister(m.RdbChangesSinceLastSave)
	registry.MustRegister(m.RdbBgsaveInProgress)
	registry.MustRegister(m.RdbLastSaveTime)
	registry.MustRegister(m.RdbLastBgsaveStatus)
	registry.MustRegister(m.RdbLastBgsaveTimeSec)
	registry.MustRegister(m.RdbCurrentBgsaveTimeSec)
	registry.MustRegister(m.RdbLastCowSize)
	registry.MustRegister(m.AofEnabled)
	registry.MustRegister(m.AofRewriteInProgress)
	registry.MustRegister(m.AofRewriteScheduled)
	registry.MustRegister(m.AofLastRewriteTimeSec)
	registry.MustRegister(m.AofCurrentRewriteTimeSec)
	registry.MustRegister(m.AofLastBgrewriteStatus)
	registry.MustRegister(m.AofLastWriteStatus)
	registry.MustRegister(m.AofLastCowSize)
}

func (m *RedisPersistenceCollector) Unregister(registry *prometheus.Registry) bool {
	if !registry.Unregister(m.Loading) {
		return false
	}
	if !registry.Unregister(m.RdbChangesSinceLastSave) {
		return false
	}
	if !registry.Unregister(m.RdbBgsaveInProgress) {
		return false
	}
	if !registry.Unregister(m.RdbLastSaveTime) {
		return false
	}
	if !registry.Unregister(m.RdbLastBgsaveStatus) {
		return false
	}
	if !registry.Unregister(m.RdbLastBgsaveTimeSec) {
		return false
	}
	if !registry.Unregister(m.RdbCurrentBgsaveTimeSec) {
		return false
	}
	if !registry.Unregister(m.RdbLastCowSize) {
		return false
	}
	if !registry.Unregister(m.AofEnabled) {
		return false
	}
	if !registry.Unregister(m.AofRewriteInProgress) {
		return false
	}
	if !registry.Unregister(m.AofRewriteScheduled) {
		return false
	}
	if !registry.Unregister(m.AofLastRewriteTimeSec) {
		return false
	}
	if !registry.Unregister(m.AofCurrentRewriteTimeSec) {
		return false
	}
	if !registry.Unregister(m.AofLastBgrewriteStatus) {
		return false
	}
	if !registry.Unregister(m.AofLastWriteStatus) {
		return false
	}
	if !registry.Unregister(m.AofLastCowSize) {
		return false
	}
	return true
}

func (m *RedisPersistenceCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	if loadingStr, ok := r["loading"]; ok {
		if loading, err := strconv.Atoi(loadingStr); err == nil {
			m.Loading.WithLabelValues(nodeName, nodeAddress).Set(float64(loading))
		}
	}
	if rdbChangersSinceLastSaveStr, ok := r["rdb_changes_since_last_save"]; ok {
		if rdbChangersSinceLastSave, err := strconv.Atoi(rdbChangersSinceLastSaveStr); err == nil {
			m.RdbChangesSinceLastSave.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbChangersSinceLastSave))
		}
	}
	if rdbBgsaveInProgressStr, ok := r["rdb_bgsave_in_progress"]; ok {
		if rdbBgsaveInProgress, err := strconv.Atoi(rdbBgsaveInProgressStr); err == nil {
			m.RdbBgsaveInProgress.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbBgsaveInProgress))
		}
	}
	if rdbLastSaveTimeStr, ok := r["rdb_last_save_time"]; ok {
		if rdbLastSaveTime, err := strconv.Atoi(rdbLastSaveTimeStr); err == nil {
			m.RdbLastSaveTime.WithLabelValues(nodeName, nodeAddress).Set(float64(rdbLastSaveTime))
		}
	}
	if rdbLastBgsaveStatus, ok := r["rdb_last_bgsave_status"]; ok {
		switch rdbLastBgsaveStatus {
		case "ok":
			m.RdbLastBgsaveStatus.WithLabelValues(nodeName, nodeAddress).Set(1)
		default:
			m.RdbLastBgsaveStatus.WithLabelValues(nodeName, nodeAddress).Set(0)
		}
	}
	if rdbLastBgsaveTimeSecStr, ok := r["rdb_last_bgsave_time_sec"]; ok {
		if rdbLastBgsaveTimeSec, err := strconv.Atoi(rdbLastBgsaveTimeSecStr); err == nil {
			m.RdbLastBgsaveTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbLastBgsaveTimeSec))
		}
	}
	if rdbCurrentBgsaveTimeSecStr, ok := r["rdb_current_bgsave_time_sec"]; ok {
		if rdbCurrentBgsaveTimeSec, err := strconv.Atoi(rdbCurrentBgsaveTimeSecStr); err == nil {
			m.RdbCurrentBgsaveTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(rdbCurrentBgsaveTimeSec))
		}
	}
	if rdbLastCowSizeStr, ok := r["rdb_last_cow_size"]; ok {
		if rdbLastCowSize, err := strconv.Atoi(rdbLastCowSizeStr); err == nil {
			m.RdbLastCowSize.WithLabelValues(nodeName, nodeAddress).Set(float64(rdbLastCowSize))
		}
	}
	if aofEnabledStr, ok := r["aof_enabled"]; ok {
		if aofEnabled, err := strconv.Atoi(aofEnabledStr); err == nil {
			m.AofEnabled.WithLabelValues(nodeName, nodeAddress).Set(float64(aofEnabled))
		}
	}
	if aofRewriteInProgressStr, ok := r["aof_rewrite_in_progress"]; ok {
		if aofRewriteInProgress, err := strconv.Atoi(aofRewriteInProgressStr); err == nil {
			m.AofRewriteInProgress.WithLabelValues(nodeName, nodeAddress).Set(
				float64(aofRewriteInProgress))
		}
	}
	if aofRewriteScheduledStr, ok := r["aof_rewrite_scheduled"]; ok {
		if aofRewriteScheduled, err := strconv.Atoi(aofRewriteScheduledStr); err == nil {
			m.AofRewriteScheduled.WithLabelValues(nodeName, nodeAddress).Set(float64(aofRewriteScheduled))
		}
	}
	if aofLastRewriteTimeSecStr, ok := r["aof_last_rewrite_time_sec"]; ok {
		if aofLastRewriteTimeSec, err := strconv.Atoi(aofLastRewriteTimeSecStr); err == nil {
			m.AofLastRewriteTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(aofLastRewriteTimeSec))
		}
	}
	if aofCurrentRewriteTimeSecStr, ok := r["aof_current_rewrite_time_sec"]; ok {
		if aofCurrentRewriteTimeSec, err := strconv.Atoi(aofCurrentRewriteTimeSecStr); err == nil {
			m.AofCurrentRewriteTimeSec.WithLabelValues(nodeName, nodeAddress).Set(
				float64(aofCurrentRewriteTimeSec))
		}
	}
	if aofLastBgrewriteStatus, ok := r["aof_last_bgrewrite_status"]; ok {
		switch aofLastBgrewriteStatus {
		case "ok":
			m.AofLastBgrewriteStatus.WithLabelValues(nodeName, nodeAddress).Set(1)
		default:
			m.AofLastBgrewriteStatus.WithLabelValues(nodeName, nodeAddress).Set(0)
		}
	}
	if aofLastWriteStatus, ok := r["aof_last_write_status"]; ok {
		switch aofLastWriteStatus {
		case "ok":
			m.AofLastWriteStatus.WithLabelValues(nodeName, nodeAddress).Set(1)
		default:
			m.AofLastWriteStatus.WithLabelValues(nodeName, nodeAddress).Set(0)
		}
	}
	if aofLastCowSizeStr, ok := r["aof_last_cow_size"]; ok {
		if aofLastCowSize, err := strconv.Atoi(aofLastCowSizeStr); err == nil {
			m.AofLastCowSize.WithLabelValues(nodeName, nodeAddress).Set(float64(aofLastCowSize))
		}
	}
	return nil
}
