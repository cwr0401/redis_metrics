package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

const (
	volatileLRU    int = iota // Evict using approximated LRU among the keys with an expire set.
	allkeysLRU                // Evict any key using approximated LRU.
	volatileLFU               // Evict using approximated LFU among the keys with an expire set.
	allkeysLFU                // Evict any key using approximated LFU.
	volatileRandom            // Remove a random key among the ones with an expire set.
	allkeysRandom             // Remove a random key, any key.
	volatileTTL               // Remove the key with the nearest expire time (minor TTL)
	noeviction                // Don't evict anything, just return an error on write operations.
)

type RedisMemoryCollector struct {
	UsedMemory             *prometheus.GaugeVec
	UsedMemoryRSS          *prometheus.GaugeVec
	UsedMemoryPeak         *prometheus.GaugeVec
	UsedMemoryOverhead     *prometheus.GaugeVec
	UsedMemoryStartup      *prometheus.GaugeVec
	UsedMemoryDataset      *prometheus.GaugeVec
	TotalSystemMemory      *prometheus.GaugeVec
	UsedMemoryLua          *prometheus.GaugeVec
	UsedMemoryScripts      *prometheus.GaugeVec
	Maxmemory              *prometheus.GaugeVec
	MaxmemoryPolicy        *prometheus.GaugeVec
	MemFragmentationRatio  *prometheus.GaugeVec
	MemAllocator           *prometheus.GaugeVec
	ActiveDefragRunning    *prometheus.GaugeVec
	LazyfreePendingObjects *prometheus.GaugeVec
	AllocatorAllocated     *prometheus.GaugeVec
	AllocatorActive        *prometheus.GaugeVec
	AllocatorResident      *prometheus.GaugeVec
	AllocatorFragRatio     *prometheus.GaugeVec
	AllocatorFragBytes     *prometheus.GaugeVec
	AllocatorRSSRatio      *prometheus.GaugeVec
	AllocatorRSSBytes      *prometheus.GaugeVec
	RSSOverheadRatio       *prometheus.GaugeVec
	RSSOverheadBytes       *prometheus.GaugeVec
	MemFragmentationBytes  *prometheus.GaugeVec
	MemNotCountedForEvict  *prometheus.GaugeVec
	MemReplicationBacklog  *prometheus.GaugeVec
	MemClientsSlaves       *prometheus.GaugeVec
	MemClientsNormal       *prometheus.GaugeVec
	MemAofBuffer           *prometheus.GaugeVec
	NumberOfCachedScripts  *prometheus.GaugeVec
}

func NewRedisMemoryCollector() *RedisMemoryCollector {

	var (
		// Memoryr
		// used_memory
		redisMemoryUsedMemory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory",
			Help:      "Total number of bytes allocated by Redis using its allocator (either standard libc, jemalloc, or an alternative allocator such as tcmalloc).",
		},
			[]string{"node_name", "node_address"})

		// used_memory_rss
		redisMemoryUsedMemoryRSS = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory_rss",
			Help:      "Number of bytes that Redis allocated as seen by the operating system (a.k.a resident set size). This is the number reported by tools such as top(1) and ps(1).",
		},
			[]string{"node_name", "node_address"})

		// used_memory_peak
		redisMemoryUsedMemoryPeak = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory_peak",
			Help:      "Peak memory consumed by Redis (in bytes), The percentage of used_memory_peak out of used_memory.",
		},
			[]string{"node_name", "node_address"})

		// used_memory_peak_perc
		//RedisMemoryUesdMemoryPeakPerc = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		//	Namespace: "redis",
		//	Subsystem: "memory",
		//	Name:      "used_memory_peak_perc",
		//	Help:      "The percentage of used_memory_peak out of used_memory.",
		//},
		//	[]string{"node_name", "node_address"})

		// used_memory_overhead
		redisMemoryUsedMemoryOverhead = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory_overhead",
			Help:      "The sum in bytes of all overheads that the server allocated for managing its internal data structures.",
		},
			[]string{"node_name", "node_address"})

		// used_memory_startup
		redisMemoryUsedMemoryStartup = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory_startup",
			Help:      "Initial amount of memory consumed by Redis at startup in bytes.",
		},
			[]string{"node_name", "node_address"})

		// used_memory_dataset
		redisMemoryUsedMemoryDataset = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory_dataset",
			Help:      "The size in bytes of the dataset (used_memory_overhead subtracted from used_memory).",
		},
			[]string{"node_name", "node_address"})

		// used_memory_dataset_perc
		//RedisMemoryUsedMemoryDatasetPerc = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		//	Namespace: "redis",
		//	Subsystem: "memory",
		//	Name:      "used_memory_dataset_perc",
		//	Help:      "The percentage of used_memory_dataset out of the net memory usage (used_memory minus used_memory_startup).",
		//},
		//	[]string{"node_name", "node_address"})

		// total_system_memory
		redisMemoryTotalSystemMemory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "total_system_memory",
			Help:      "The total amount of memory that the Redis host has.",
		},
			[]string{"node_name", "node_address"})

		// used_memory_lua
		redisMemoryUsedMemoryLua = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory_lua",
			Help:      "Number of bytes used by the Lua engine.",
		},
			[]string{"node_name", "node_address"})

		// used_memory_scripts
		redisMemoryUsedMemoryScripts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "used_memory_scripts",
			Help:      "Number of bytes used by cached Lua scripts.",
		},
			[]string{"node_name", "node_address"})

		// maxmemory
		redisMemoryMaxmemory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "maxmemory",
			Help:      "The value of the maxmemory configuration directive.",
		},
			[]string{"node_name", "node_address"})

		// maxmemory_policy
		//
		redisMemoryMaxmemoryPolicy = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "maxmemory_policy",
			Help:      "The value of the maxmemory-policy configuration directive",
		},
			[]string{"node_name", "node_address"})

		// mem_fragmentation_ratio
		redisMemoryMemFragmentationRatio = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_fragmentation_ratio",
			Help:      "Ratio between used_memory_rss and used_memory.",
		},
			[]string{"node_name", "node_address"})

		// mem_allocator
		redisMemoryMemAllocator = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_allocator",
			Help:      "Memory allocator, chosen at compile time.",
		},
			[]string{"node_name", "node_address", "allocator"})

		// active_defrag_running
		redisMemoryActiveDefragRunning = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "active_defrag_running",
			Help:      "Flag indicating if active defragmentation is active.",
		},
			[]string{"node_name", "node_address"})

		// lazyfree_pending_objects
		redisMemoryLazyfreePendingObjects = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "lazyfree_pending_objects",
			Help:      "The number of objects waiting to be freed (as a result of calling UNLINK, or FLUSHDB and FLUSHALL with the ASYNC option).",
		},
			[]string{"node_name", "node_address"})

		// allocator_allocated
		redisMemoryAllocatorAllocated = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "allocator_allocated",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// allocator_active
		redisMemoryAllocatorActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "allocator_active",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// allocator_resident
		redisMemoryAllocatorResident = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "allocator_resident",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// allocator_frag_ratio
		redisMemoryAllocatorFragRatio = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "allocator_frag_ratio",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// allocator_frag_bytes
		redisMemoryAllocatorFragBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "allocator_frag_bytes",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// allocator_rss_ratio
		redisMemoryAllocatorRSSRatio = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "allocator_rss_ratio",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// allocator_rss_bytes
		redisMemoryAllocatorRSSBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "allocator_rss_bytes",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// rss_overhead_ratio
		redisMemoryRSSOverheadRatio = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "rss_overhead_ratio",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// rss_overhead_bytes
		redisMemoryRSSOverheadBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "rss_overhead_bytes",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// mem_fragmentation_bytes
		redisMemoryMemFragmentationBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_fragmentation_bytes",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// mem_not_counted_for_evict
		redisMemoryMemNotCountedForEvict = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_not_counted_for_evict",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// mem_replication_backlog
		redisMemoryMemReplicationBacklog = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_replication_backlog",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// mem_clients_slaves
		redisMemoryMemClientsSlaves = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_clients_slaves",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// mem_clients_normal
		redisMemoryMemClientsNormal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_clients_normal",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// mem_aof_buffer
		redisMemoryMemAofBuffer = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "mem_aof_buffer",
			Help:      "",
		},
			[]string{"node_name", "node_address"})

		// number_of_cached_scripts
		redisMemoryNumberOfCachedScripts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "memory",
			Name:      "number_of_cached_scripts",
			Help:      "",
		},
			[]string{"node_name", "node_address"})
	)
	return &RedisMemoryCollector{
		redisMemoryUsedMemory,
		redisMemoryUsedMemoryRSS,
		redisMemoryUsedMemoryPeak,
		redisMemoryUsedMemoryOverhead,
		redisMemoryUsedMemoryStartup,
		redisMemoryUsedMemoryDataset,
		redisMemoryTotalSystemMemory,
		redisMemoryUsedMemoryLua,
		redisMemoryUsedMemoryScripts,
		redisMemoryMaxmemory,
		redisMemoryMaxmemoryPolicy,
		redisMemoryMemFragmentationRatio,
		redisMemoryMemAllocator,
		redisMemoryActiveDefragRunning,
		redisMemoryLazyfreePendingObjects,
		redisMemoryAllocatorAllocated,
		redisMemoryAllocatorActive,
		redisMemoryAllocatorResident,
		redisMemoryAllocatorFragRatio,
		redisMemoryAllocatorFragBytes,
		redisMemoryAllocatorRSSRatio,
		redisMemoryAllocatorRSSBytes,
		redisMemoryRSSOverheadRatio,
		redisMemoryRSSOverheadBytes,
		redisMemoryMemFragmentationBytes,
		redisMemoryMemNotCountedForEvict,
		redisMemoryMemReplicationBacklog,
		redisMemoryMemClientsSlaves,
		redisMemoryMemClientsNormal,
		redisMemoryMemAofBuffer,
		redisMemoryNumberOfCachedScripts,
	}
}

func (m *RedisMemoryCollector) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(m.UsedMemory)
	registry.MustRegister(m.UsedMemoryRSS)
	registry.MustRegister(m.UsedMemoryPeak)
	// registry.MustRegister(m.UesdMemoryPeakPerc)
	registry.MustRegister(m.UsedMemoryOverhead)
	registry.MustRegister(m.UsedMemoryStartup)
	registry.MustRegister(m.UsedMemoryDataset)
	// registry.MustRegister(m.UsedMemoryDatasetPerc)
	registry.MustRegister(m.TotalSystemMemory)
	registry.MustRegister(m.UsedMemoryLua)
	registry.MustRegister(m.UsedMemoryScripts)
	registry.MustRegister(m.Maxmemory)
	registry.MustRegister(m.MaxmemoryPolicy)
	registry.MustRegister(m.MemFragmentationRatio)
	registry.MustRegister(m.MemAllocator)
	registry.MustRegister(m.ActiveDefragRunning)
	registry.MustRegister(m.LazyfreePendingObjects)
	registry.MustRegister(m.AllocatorAllocated)
	registry.MustRegister(m.AllocatorActive)
	registry.MustRegister(m.AllocatorResident)
	registry.MustRegister(m.AllocatorFragRatio)
	registry.MustRegister(m.AllocatorFragBytes)
	registry.MustRegister(m.AllocatorRSSRatio)
	registry.MustRegister(m.AllocatorRSSBytes)
	registry.MustRegister(m.RSSOverheadRatio)
	registry.MustRegister(m.RSSOverheadBytes)
	registry.MustRegister(m.MemFragmentationBytes)
	registry.MustRegister(m.MemNotCountedForEvict)
	registry.MustRegister(m.MemReplicationBacklog)
	registry.MustRegister(m.MemClientsSlaves)
	registry.MustRegister(m.MemClientsNormal)
	registry.MustRegister(m.MemAofBuffer)
	registry.MustRegister(m.NumberOfCachedScripts)
}

func (m *RedisMemoryCollector) Unregister(registry *prometheus.Registry) bool {
	if !registry.Unregister(m.UsedMemory) {
		return false
	}
	if !registry.Unregister(m.UsedMemoryRSS) {
		return false
	}
	if !registry.Unregister(m.UsedMemoryPeak) {
		return false
	}
	// registry.MustRegister(m.UesdMemoryPeakPerc)
	if !registry.Unregister(m.UsedMemoryOverhead) {
		return false
	}
	if !registry.Unregister(m.UsedMemoryStartup) {
		return false
	}
	if !registry.Unregister(m.UsedMemoryDataset) {
		return false
	}
	// registry.MustRegister(m.UsedMemoryDatasetPerc)
	if !registry.Unregister(m.TotalSystemMemory) {
		return false
	}
	if !registry.Unregister(m.UsedMemoryLua) {
		return false
	}
	if !registry.Unregister(m.UsedMemoryScripts) {
		return false
	}
	if !registry.Unregister(m.Maxmemory) {
		return false
	}
	if !registry.Unregister(m.MaxmemoryPolicy) {
		return false
	}
	if !registry.Unregister(m.MemFragmentationRatio) {
		return false
	}
	if !registry.Unregister(m.MemAllocator) {
		return false
	}
	if !registry.Unregister(m.ActiveDefragRunning) {
		return false
	}
	if !registry.Unregister(m.LazyfreePendingObjects) {
		return false
	}
	if !registry.Unregister(m.AllocatorAllocated) {
		return false
	}
	if !registry.Unregister(m.AllocatorActive) {
		return false
	}
	if !registry.Unregister(m.AllocatorResident) {
		return false
	}
	if !registry.Unregister(m.AllocatorFragRatio) {
		return false
	}
	if !registry.Unregister(m.AllocatorFragBytes) {
		return false
	}
	if !registry.Unregister(m.AllocatorRSSRatio) {
		return false
	}
	if !registry.Unregister(m.AllocatorRSSBytes) {
		return false
	}
	if !registry.Unregister(m.RSSOverheadRatio) {
		return false
	}
	if !registry.Unregister(m.RSSOverheadBytes) {
		return false
	}
	if !registry.Unregister(m.MemFragmentationBytes) {
		return false
	}
	if !registry.Unregister(m.MemNotCountedForEvict) {
		return false
	}
	if !registry.Unregister(m.MemReplicationBacklog) {
		return false
	}
	if !registry.Unregister(m.MemClientsSlaves) {
		return false
	}
	if !registry.Unregister(m.MemClientsNormal) {
		return false
	}
	if !registry.Unregister(m.MemAofBuffer) {
		return false
	}
	if !registry.Unregister(m.NumberOfCachedScripts) {
		return false
	}
	return true
}

func (m *RedisMemoryCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	if usedMemoryStr, ok := r["used_memory"]; ok {
		if usedMemory, err := strconv.Atoi(usedMemoryStr); err == nil {
			m.UsedMemory.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemory))
		}
	}
	if usedMemoryRSSStr, ok := r["used_memory_rss"]; ok {
		if usedMemoryRSS, err := strconv.Atoi(usedMemoryRSSStr); err == nil {
			m.UsedMemoryRSS.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemoryRSS))
		}
	}
	if usedMemoryPeakStr, ok := r["used_memory_peak"]; ok {
		if usedMemoryPeak, err := strconv.Atoi(usedMemoryPeakStr); err == nil {
			m.UsedMemoryPeak.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemoryPeak))
		}
	}

	//usedMemoryPeakPerc := "0%"
	//if usedMemoryPeakPerc, ok = r["used_memory_peak_perc"]; !ok {
	//	usedMemoryPeakPerc = "0%"
	//}

	if usedMemoryOverheadStr, ok := r["used_memory_overhead"]; ok {
		if usedMemoryOverhead, err := strconv.Atoi(usedMemoryOverheadStr); err == nil {
			m.UsedMemoryOverhead.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemoryOverhead))
		}
	}
	if usedMemoryStratupStr, ok := r["used_memory_startup"]; ok {
		if usedMemoryStratup, err := strconv.Atoi(usedMemoryStratupStr); err == nil {
			m.UsedMemoryStartup.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemoryStratup))
		}
	}
	if usedMemoryDatasetStr, ok := r["used_memory_dataset"]; ok {
		if usedMemoryDataset, err := strconv.Atoi(usedMemoryDatasetStr); err == nil {
			m.UsedMemoryDataset.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemoryDataset))
		}
	}
	//if usedMemoryDatasetPerc, ok := r["used_memory_dataset_perc"]; ok {
	//}
	if allocatorAllocatedStr, ok := r["allocator_allocated"]; ok {
		if allocatorAllocated, err := strconv.Atoi(allocatorAllocatedStr); err == nil {
			m.AllocatorAllocated.WithLabelValues(nodeName, nodeAddress).Set(float64(allocatorAllocated))
		}
	}
	if allocatorActiveStr, ok := r["allocator_active"]; ok {
		if allocatorActive, err := strconv.Atoi(allocatorActiveStr); err == nil {
			m.AllocatorActive.WithLabelValues(nodeName, nodeAddress).Set(float64(allocatorActive))
		}
	}
	if allocatorResidentStr, ok := r["allocator_resident"]; ok {
		if allocatorResident, err := strconv.Atoi(allocatorResidentStr); err == nil {
			m.AllocatorResident.WithLabelValues(nodeName, nodeAddress).Set(float64(allocatorResident))
		}
	}
	if totalSystemMemoryStr, ok := r["total_system_memory"]; ok {
		if totalSystemMemory, err := strconv.Atoi(totalSystemMemoryStr); err == nil {
			m.TotalSystemMemory.WithLabelValues(nodeName, nodeAddress).Set(float64(totalSystemMemory))
		}
	}
	if usedMemoryLuaStr, ok := r["used_memory_lua"]; ok {
		if usedMemoryLua, err := strconv.Atoi(usedMemoryLuaStr); err == nil {
			m.UsedMemoryLua.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemoryLua))
		}
	}
	if usedMemoryScriptsStr, ok := r["used_memory_scripts"]; ok {
		if usedMemoryScripts, err := strconv.Atoi(usedMemoryScriptsStr); err == nil {
			m.UsedMemoryScripts.WithLabelValues(nodeName, nodeAddress).Set(float64(usedMemoryScripts))
		}
	}
	if numberOfCachedScriptsStr, ok := r["number_of_cached_scripts"]; ok {
		if numberOfCachedScripts, err := strconv.Atoi(numberOfCachedScriptsStr); err == nil {
			m.NumberOfCachedScripts.WithLabelValues(nodeName, nodeAddress).Set(float64(numberOfCachedScripts))
		}
	}
	if maxmemoryStr, ok := r["maxmemory"]; ok {
		if maxmemory, err := strconv.Atoi(maxmemoryStr); err == nil {
			m.Maxmemory.WithLabelValues(nodeName, nodeAddress).Set(float64(maxmemory))
		}
	}
	if maxmemoryPolicyStr, ok := r["maxmemory_policy"]; ok {
		maxmemoryPolicy := noeviction
		switch maxmemoryPolicyStr {
		case "volatile-lru":
			maxmemoryPolicy = volatileLRU
		case "allkeys-lru":
			maxmemoryPolicy = allkeysLRU
		case "volatile-lfu":
			maxmemoryPolicy = volatileLFU
		case "allkeys-lfu":
			maxmemoryPolicy = allkeysLFU
		case "volatile-random":
			maxmemoryPolicy = volatileRandom
		case "allkeys-random":
			maxmemoryPolicy = allkeysRandom
		case "volatile-ttl":
			maxmemoryPolicy = volatileTTL
		case "noeviction":
			maxmemoryPolicy = noeviction
		}
		m.MaxmemoryPolicy.WithLabelValues(nodeName, nodeAddress).Set(float64(maxmemoryPolicy))
	}
	if allocatorFragRatioStr, ok := r["allocator_frag_ratio"]; ok {
		if allocatorFragRatio, err := strconv.ParseFloat(allocatorFragRatioStr, 64); err == nil {
			m.AllocatorFragRatio.WithLabelValues(nodeName, nodeAddress).Set(allocatorFragRatio)
		}
	}
	if allocatorFragBytesStr, ok := r["allocator_frag_bytes"]; ok {
		if allocatorFragBytes, err := strconv.Atoi(allocatorFragBytesStr); err == nil {
			m.AllocatorFragBytes.WithLabelValues(nodeName, nodeAddress).Set(float64(allocatorFragBytes))
		}
	}
	if rssOverheadRatioStr, ok := r["rss_overhead_ratio"]; ok {
		if rssOverheadRatio, err := strconv.ParseFloat(rssOverheadRatioStr, 64); err == nil {
			m.RSSOverheadRatio.WithLabelValues(nodeName, nodeAddress).Set(rssOverheadRatio)
		}
	}
	if rssOverheadBytesStr, ok := r["rss_overhead_bytes"]; ok {
		if rssOverheadBytes, err := strconv.Atoi(rssOverheadBytesStr); err == nil {
			m.RSSOverheadBytes.WithLabelValues(nodeName, nodeAddress).Set(float64(rssOverheadBytes))
		}
	}
	if memFragmentationRatioStr, ok := r["mem_fragmentation_ratio"]; ok {
		if memFragmentationRatio, err := strconv.ParseFloat(memFragmentationRatioStr, 64); err == nil {
			m.MemFragmentationRatio.WithLabelValues(nodeName, nodeAddress).Set(memFragmentationRatio)
		}
	}
	if memFragmentationBytesStr, ok := r["mem_fragmentation_bytes"]; ok {
		if memFragmentationBytes, err := strconv.Atoi(memFragmentationBytesStr); err == nil {
			m.MemFragmentationBytes.WithLabelValues(nodeName, nodeAddress).Set(float64(memFragmentationBytes))
		}
	}
	if memNotCountedForEvictStr, ok := r["mem_not_counted_for_evict"]; ok {
		if memNotCountedForEvict, err := strconv.Atoi(memNotCountedForEvictStr); err == nil {
			m.MemNotCountedForEvict.WithLabelValues(nodeName, nodeAddress).Set(float64(memNotCountedForEvict))
		}
	}
	if memReplicationBacklogStr, ok := r["mem_replication_backlog"]; ok {
		if memReplicationBacklog, err := strconv.Atoi(memReplicationBacklogStr); err == nil {
			m.MemReplicationBacklog.WithLabelValues(nodeName, nodeAddress).Set(float64(memReplicationBacklog))
		}
	}
	if memClientsSlavesStr, ok := r["mem_clients_slaves"]; ok {
		if memClientsSlaves, err := strconv.Atoi(memClientsSlavesStr); err == nil {
			m.MemClientsSlaves.WithLabelValues(nodeName, nodeAddress).Set(float64(memClientsSlaves))
		}
	}
	if memClientsNormalStr, ok := r["mem_clients_normal"]; ok {
		if memClientsNormal, err := strconv.Atoi(memClientsNormalStr); err == nil {
			m.MemClientsNormal.WithLabelValues(nodeName, nodeAddress).Set(float64(memClientsNormal))
		}
	}
	if memAofBufferStr, ok := r["mem_aof_buffer"]; ok {
		if memAofBuffer, err := strconv.Atoi(memAofBufferStr); err == nil {
			m.MemAofBuffer.WithLabelValues(nodeName, nodeAddress).Set(float64(memAofBuffer))
		}
	}
	if memAllocator, ok := r["mem_allocator"]; ok {
		m.MemAllocator.WithLabelValues(nodeName, nodeAddress, memAllocator).Set(1)
	}
	if activeDefragRunningStr, ok := r["active_defrag_running"]; ok {
		if activeDefragRunning, err := strconv.Atoi(activeDefragRunningStr); err == nil {
			m.ActiveDefragRunning.WithLabelValues(nodeName, nodeAddress).Set(float64(activeDefragRunning))
		}
	}
	if lazyfreePendingObjectsStr, ok := r["lazyfree_pending_objects"]; ok {
		if lazyfreePendingObjects, err := strconv.Atoi(lazyfreePendingObjectsStr); err == nil {
			m.LazyfreePendingObjects.WithLabelValues(nodeName, nodeAddress).Set(float64(lazyfreePendingObjects))
		}
	}
	return nil
}
