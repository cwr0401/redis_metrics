package metrics

import (
	"time"

	"github.com/cwr0401/redis_metrics/config"
	"github.com/cwr0401/redis_metrics/metrics/info"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry                        = prometheus.NewRegistry()
	gather      prometheus.Gatherer = registry
	promHandler                     = promhttp.HandlerFor(gather, promhttp.HandlerOpts{})
	Handler                         = promhttp.InstrumentMetricHandler(registry, promHandler)
)

func init() {
	// log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	// log.SetLevel(log.DebugLevel)

	// Metrics have to be registered to be exposed:
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	// Server
	registry.MustRegister(info.RedisServerUp)
	registry.MustRegister(info.RedisServerInfo)
	registry.MustRegister(info.RedisServerUptimeInSeconds)
	registry.MustRegister(info.RedisServerUptimeInDays)
	registry.MustRegister(info.RedisServerHz)
	registry.MustRegister(info.RedisServerConfiguredHz)
	registry.MustRegister(info.RedisServerLruClock)

	// Clients
	registry.MustRegister(info.RedisClientsConnectedClients)
	registry.MustRegister(info.RedisClientsClientRecentMaxInputBuffer)
	registry.MustRegister(info.RedisClientsClientRecentMaxOutputBuffer)
	//registry.MustRegister(redis.RedisClientsClientRecentMaxInputBuffer)
	//registry.MustRegister(redis.RedisClientsClientRecentMaxOutputBuffer)
	registry.MustRegister(info.RedisClientsBlockedClients)

	// Memory
	registry.MustRegister(info.RedisMemoryUsedMemory)
	registry.MustRegister(info.RedisMemoryUsedMemoryRSS)
	registry.MustRegister(info.RedisMemoryUsedMemoryPeak)
	// registry.MustRegister(info.RedisMemoryUesdMemoryPeakPerc)
	registry.MustRegister(info.RedisMemoryUsedMemoryOverhead)
	registry.MustRegister(info.RedisMemoryUsedMemoryStartup)
	registry.MustRegister(info.RedisMemoryUsedMemoryDataset)
	// registry.MustRegister(info.RedisMemoryUsedMemoryDatasetPerc)
	registry.MustRegister(info.RedisMemoryTotalSystemMemory)
	registry.MustRegister(info.RedisMemoryUsedMemoryLua)
	registry.MustRegister(info.RedisMemoryUsedMemoryScripts)
	registry.MustRegister(info.RedisMemoryMaxmemory)
	registry.MustRegister(info.RedisMemoryMaxmemoryPolicy)
	registry.MustRegister(info.RedisMemoryMemFragmentationRatio)
	registry.MustRegister(info.RedisMemoryMemAllocator)
	registry.MustRegister(info.RedisMemoryActiveDefragRunning)
	registry.MustRegister(info.RedisMemoryLazyfreePendingObjects)
	registry.MustRegister(info.RedisMemoryAllocatorAllocated)
	registry.MustRegister(info.RedisMemoryAllocatorActive)
	registry.MustRegister(info.RedisMemoryAllocatorResident)
	registry.MustRegister(info.RedisMemoryAllocatorFragRatio)
	registry.MustRegister(info.RedisMemoryAllocatorFragBytes)
	registry.MustRegister(info.RedisMemoryAllocatorRSSRatio)
	registry.MustRegister(info.RedisMemoryAllocatorRSSBytes)
	registry.MustRegister(info.RedisMemoryRSSOverheadRatio)
	registry.MustRegister(info.RedisMemoryRSSOverheadBytes)
	registry.MustRegister(info.RedisMemoryMemFragmentationBytes)
	registry.MustRegister(info.RedisMemoryMemNotCountedForEvict)
	registry.MustRegister(info.RedisMemoryMemReplicationBacklog)
	registry.MustRegister(info.RedisMemoryMemClientsSlaves)
	registry.MustRegister(info.RedisMemoryMemClientsNormal)
	registry.MustRegister(info.RedisMemoryMemAofBuffer)
	registry.MustRegister(info.RedisMemoryNumberOfCachedScripts)

	// Persistence
	registry.MustRegister(info.RedisPersistenceLoading)
	registry.MustRegister(info.RedisPersistenceRdbChangesSinceLastSave)
	registry.MustRegister(info.RedisPersistenceRdbBgsaveInProgress)
	registry.MustRegister(info.RedisPersistenceRdbLastSaveTime)
	registry.MustRegister(info.RedisPersistenceRdbLastBgsaveStatus)
	registry.MustRegister(info.RedisPersistenceRdbLastBgsaveTimeSec)
	registry.MustRegister(info.RedisPersistenceRdbCurrentBgsaveTimeSec)
	registry.MustRegister(info.RedisPersistenceRdbLastCowSize)
	registry.MustRegister(info.RedisPersistenceAofEnabled)
	registry.MustRegister(info.RedisPersistenceAofRewriteInProgress)
	registry.MustRegister(info.RedisPersistenceAofRewriteScheduled)
	registry.MustRegister(info.RedisPersistenceAofLastRewriteTimeSec)
	registry.MustRegister(info.RedisPersistenceAofCurrentRewriteTimeSec)
	registry.MustRegister(info.RedisPersistenceAofLastBgrewriteStatus)
	registry.MustRegister(info.RedisPersistenceAofLastWriteStatus)
	registry.MustRegister(info.RedisPersistenceAofLastCowSize)

	// Stats
	registry.MustRegister(info.RedisStatsTotalConnectionsReceived)
	registry.MustRegister(info.RedisStatsTotalCommandsProcessed)
	registry.MustRegister(info.RedisStatsInstantaneousOpsPerSec)
	registry.MustRegister(info.RedisStatsTotalNetInputBytes)
	registry.MustRegister(info.RedisStatsTotalNetOutputBytes)
	registry.MustRegister(info.RedisStatsInstantaneousInputKbps)
	registry.MustRegister(info.RedisStatsInstantaneousOutputKbps)
	registry.MustRegister(info.RedisStatsRejectedConnections)
	registry.MustRegister(info.RedisStatsSyncFull)
	registry.MustRegister(info.RedisStatsSyncPartialOK)
	registry.MustRegister(info.RedisStatsSyncPartialErr)
	registry.MustRegister(info.RedisStatsExpiredKeys)
	registry.MustRegister(info.RedisStatsExpiredStalePerc)
	registry.MustRegister(info.RedisStatsExpiredTimeCapReachedCount)
	registry.MustRegister(info.RedisStatsEvictedKeys)
	registry.MustRegister(info.RedisStatsKeyspaceHits)
	registry.MustRegister(info.RedisStatsKeyspaceMisses)
	registry.MustRegister(info.RedisStatsPubsubChannels)
	registry.MustRegister(info.RedisStatsPubsubPatterns)
	registry.MustRegister(info.RedisStatsLatestForkUsec)
	registry.MustRegister(info.RedisStatsMigrateCachedSockets)
	registry.MustRegister(info.RedisStatsSlaveExpiresTrackedKeys)
	registry.MustRegister(info.RedisStatsActiveDefragHits)
	registry.MustRegister(info.RedisStatsActiveDefragMisses)
	registry.MustRegister(info.RedisStatsActiveDefragKeyHits)
	registry.MustRegister(info.RedisStatsActiveDefragKeyMisses)

	// Replication
	registry.MustRegister(info.RedisReplicationRole)
	registry.MustRegister(info.RedisReplicationConnectedSlaves)
	//registry.MustRegister(info.RedisReplicationMasterRepl)
	registry.MustRegister(info.RedisReplicationMasterReplOffset)
	registry.MustRegister(info.RedisReplicationSecondReplOffset)
	registry.MustRegister(info.RedisReplicationReplBacklogActive)
	registry.MustRegister(info.RedisReplicationReplBacklogSize)
	registry.MustRegister(info.RedisReplicationReplBacklogFirstByteOffset)
	registry.MustRegister(info.RedisReplicationReplBacklogHistlen)
	registry.MustRegister(info.RedisReplicationMasterHostPort)
	registry.MustRegister(info.RedisReplicationMasterLinkStatus)
	registry.MustRegister(info.RedisReplicationMasterLastIOSecondsAgo)
	registry.MustRegister(info.RedisReplicationMasterSyncInProgress)
	registry.MustRegister(info.RedisReplicationSlaveReplOffset)
	registry.MustRegister(info.RedisReplicationSlavePriority)
	registry.MustRegister(info.RedisReplicationSlaveReadOnly)
	registry.MustRegister(info.RedisReplicationMasterSyncLeftBytes)
	registry.MustRegister(info.RedisReplicationMasterSyncLastIOSecondsAgo)
	registry.MustRegister(info.RedisReplicationMasterLinkDownSinceSeconds)
	registry.MustRegister(info.RedisReplicationMinSlavesGoodSlaves)
	registry.MustRegister(info.RedisReplicationSlaveState)
	registry.MustRegister(info.RedisReplicationSlaveOffset)
	registry.MustRegister(info.RedisReplicationSlaveLag)

	// CPU
	registry.MustRegister(info.RedisCPUUsedCpuSys)
	registry.MustRegister(info.RedisCPUUsedCpuUser)
	registry.MustRegister(info.RedisCPUUsedCpuSysChildren)
	registry.MustRegister(info.RedisCPUUsedCpuUserChildren)

	// Keyspace
	registry.MustRegister(info.RedisKeyspaceDBKeys)
	registry.MustRegister(info.RedisKeyspaceDBExpires)
	registry.MustRegister(info.RedisKeyspaceDBAvgTTL)

	// Sentinel
	registry.MustRegister(info.RedisSentinelUp)
	registry.MustRegister(info.RedisSentinelSentinelMasters)
	registry.MustRegister(info.RedisSentinelSentinelTilt)
	registry.MustRegister(info.RedisSentinelSentinelRunningScripts)
	registry.MustRegister(info.RedisSentinelSentinelScriptsQueueLength)
	registry.MustRegister(info.RedisSentinelSentinelSimulateFailureFlags)
	registry.MustRegister(info.RedisSentinelMaster)
	registry.MustRegister(info.RedisSentinelMastersInfo)
	registry.MustRegister(info.RedisSentinelSlavesInfo)
	registry.MustRegister(info.RedisSentinelSentinelsInfo)

	// Cluster
	registry.MustRegister(info.RedisClusterClusterEnabled)
}

func RedisCollector(c *config.RedisConfig, duration time.Duration) error {
	for _, standalone := range c.Standalone {
		standaloneOptions := standalone.RedisOptions()
		go RedisServerMetrics(standalone.Name, standaloneOptions, duration)
	}
	for _, sentinel := range c.Sentinel {
		sentinelOptions := sentinel.RedisOptions()
		go RedisSentinelMetrics(
			sentinel.Name, sentinelOptions, duration,
			sentinel.DiscoveryMaster, sentinel.DiscoverySlave, sentinel.DiscoverySentinel)
	}
	//for _, cluster := range c.Cluster {
	//	clusterOptions := cluster.RedisClusterOptions()
	//	fmt.Println(clusterOptions)
	//}

	for {
		time.Sleep(time.Minute * 60)
	}

	return nil
}
