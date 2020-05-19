package info

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	RedisReplicationRoleMaster int = iota
	RedisReplicationRoleSlave
)

const (
	RedisReplicationMasterLinkUp int = iota
	RedisReplicationMasterLinkDown
)

const (
	RedisReplicationSlaveStateOnline int = iota
	RedisReplicationSlaveStateOffline
)

type RedisReplicationCollector struct {
	Role                       *prometheus.GaugeVec
	ConnectedSlaves            *prometheus.GaugeVec
	MasterReplOffset           *prometheus.GaugeVec
	SecondReplOffset           *prometheus.GaugeVec
	ReplBacklogActive          *prometheus.GaugeVec
	ReplBacklogSize            *prometheus.GaugeVec
	ReplBacklogFirstByteOffset *prometheus.GaugeVec
	ReplBacklogHistlen         *prometheus.GaugeVec
	MasterHostPort             *prometheus.CounterVec
	MasterLinkStatus           *prometheus.GaugeVec
	MasterLastIOSecondsAgo     *prometheus.GaugeVec
	MasterSyncInProgress       *prometheus.GaugeVec
	SlaveReplOffset            *prometheus.GaugeVec
	SlavePriority              *prometheus.GaugeVec
	SlaveReadOnly              *prometheus.GaugeVec
	MasterSyncLeftBytes        *prometheus.GaugeVec
	MasterSyncLastIOSecondsAgo *prometheus.GaugeVec
	MasterLinkDownSinceSeconds *prometheus.GaugeVec
	MinSlavesGoodSlaves        *prometheus.GaugeVec
	SlaveState                 *prometheus.GaugeVec
	SlaveOffset                *prometheus.GaugeVec
	SlaveLag                   *prometheus.GaugeVec
}

func NewRedisReplicationCollector() *RedisReplicationCollector {
	var (
		// Replication

		// role:master
		//connected_slaves:2
		//slave0:ip=172.25.1.5,port=6379,state=online,offset=1137408,lag=0
		//slave1:ip=172.25.1.4,port=6379,state=online,offset=1137408,lag=0
		//master_replid:fd3d751298a80be53bb6314890677ba7117244f2
		//master_replid2:0000000000000000000000000000000000000000
		//master_repl_offset:1137408
		//second_repl_offset:-1
		//repl_backlog_active:1
		//repl_backlog_size:1048576
		//repl_backlog_first_byte_offset:88833
		//repl_backlog_histlen:1048576

		// role:slave
		//master_host:172.25.1.3
		//master_port:6379
		//master_link_status:up
		//master_last_io_seconds_ago:0
		//master_sync_in_progress:0
		//slave_repl_offset:1143525
		//slave_priority:100
		//slave_read_only:1
		//connected_slaves:0
		//master_replid:fd3d751298a80be53bb6314890677ba7117244f2
		//master_replid2:0000000000000000000000000000000000000000
		//master_repl_offset:1143525
		//second_repl_offset:-1
		//repl_backlog_active:1
		//repl_backlog_size:1048576
		//repl_backlog_first_byte_offset:94950
		//repl_backlog_histlen:1048576

		// role:master
		redisReplicationRole = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "role",
			Help:      "Value is Master(0) if the instance is replica of no one, or Slave(1) if the instance is a replica of some master instance. Note that a replica can be master of another replica (chained replication). ",
		},
			[]string{"node_name", "node_address"})

		// connected_slaves:0
		redisReplicationConnectedSlaves = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "connected_slaves",
			Help:      "Number of connected replicas.",
		},
			[]string{"node_name", "node_address"})

		// master_replid:ec1615c178d40a7bed769a9753c2b15c42bd0ba9
		// master_replid2:0000000000000000000000000000000000000000
		//RedisReplicationMasterRepl = prometheus.NewCounterVec(prometheus.CounterOpts{
		//	Namespace: "redis",
		//	Subsystem: "replication",
		//	Name: "master_repl",
		//	Help: "The replication ID of the Redis server. The secondary replication ID, used for PSYNC after a failover.",
		//},
		//	[]string{"node_name", "node_address", "id", "id2"})

		//RedisReplicationMasterReplid2 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		//	Namespace: "redis",
		//	Subsystem: "replication",
		//	Name: "master_replid2",
		//	Help: "The secondary replication ID, used for PSYNC after a failover.",
		//},
		//	[]string{})

		// master_repl_offset:0
		redisReplicationMasterReplOffset = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master_repl_offset",
			Help:      "The server's current replication offset.",
		},
			[]string{"node_name", "node_address"})

		// second_repl_offset:-1
		redisReplicationSecondReplOffset = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "second_repl_offset",
			Help:      "The offset up to which replication IDs are accepted.",
		},
			[]string{"node_name", "node_address"})

		// repl_backlog_active:0
		redisReplicationReplBacklogActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "repl_backlog_active",
			Help:      "Flag indicating replication backlog is active.",
		},
			[]string{"node_name", "node_address"})

		// repl_backlog_size:1048576
		redisReplicationReplBacklogSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "repl_backlog_size",
			Help:      "Total size in bytes of the replication backlog buffer.",
		},
			[]string{"node_name", "node_address"})

		// repl_backlog_first_byte_offset:0
		redisReplicationReplBacklogFirstByteOffset = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "repl_backlog_first_byte_offset",
			Help:      "The master offset of the replication backlog buffer.",
		},
			[]string{"node_name", "node_address"})

		// repl_backlog_histlen:0
		redisReplicationReplBacklogHistlen = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "repl_backlog_histlen",
			Help:      "Size in bytes of the data in the replication backlog buffer.",
		},
			[]string{"node_name", "node_address"})

		// master_host
		redisReplicationMasterHostPort = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master",
			Help:      "Host or IP address of the master. Master listening TCP port.",
		},
			[]string{"node_name", "node_address", "host", "port"})

		// master_port
		//RedisReplicationMasterPort = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		//	Namespace: "redis",
		//	Subsystem: "replication",
		//	Name: "master_port",
		//	Help: "Master listening TCP port.",
		//},
		//	[]string{})

		// master_link_status
		redisReplicationMasterLinkStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master_link_status",
			Help:      "Status of the link (up/down).",
		},
			[]string{"node_name", "node_address"})

		// master_last_io_seconds_ago
		redisReplicationMasterLastIOSecondsAgo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master_last_io_seconds_ago",
			Help:      "Number of seconds since the last interaction with master.",
		},
			[]string{"node_name", "node_address"})

		// master_sync_in_progress
		redisReplicationMasterSyncInProgress = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master_sync_in_progress",
			Help:      "Indicate the master is syncing to the replica.",
		},
			[]string{"node_name", "node_address"})

		// slave_repl_offset
		redisReplicationSlaveReplOffset = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "slave_repl_offset",
			Help:      "The replication offset of the replica instance.",
		},
			[]string{"node_name", "node_address"})

		// slave_priority
		redisReplicationSlavePriority = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "slave_priority",
			Help:      "The priority of the instance as a candidate for failover",
		},
			[]string{"node_name", "node_address"})

		// slave_read_only
		redisReplicationSlaveReadOnly = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "slave_read_only",
			Help:      "Flag indicating if the replica is read-only.",
		},
			[]string{"node_name", "node_address"})

		// master_sync_left_bytes
		redisReplicationMasterSyncLeftBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master_sync_left_bytes",
			Help:      "Number of bytes left before syncing is complete.",
		},
			[]string{"node_name", "node_address"})

		// master_sync_last_io_seconds_ago
		redisReplicationMasterSyncLastIOSecondsAgo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master_sync_last_io_seconds_ago",
			Help:      "Number of seconds since last transfer I/O during a SYNC operation.",
		},
			[]string{"node_name", "node_address"})

		// master_link_down_since_seconds
		redisReplicationMasterLinkDownSinceSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "master_link_down_since_seconds",
			Help:      "Number of seconds since the link is down.",
		},
			[]string{"node_name", "node_address"})

		// min_slaves_good_slaves
		redisReplicationMinSlavesGoodSlaves = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "min_slaves_good_slaves",
			Help:      "Number of replicas currently considered good.",
		},
			[]string{"node_name", "node_address"})

		// slaveXXX
		redisReplicationSlaveState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "slave_state",
			Help:      "Slave id, IP address, port, state(Value 0 is online, Value 1 is offset)",
		},
			[]string{
				"node_name",
				"node_address",
				"slave_id",
				"slave_addr",
			})
		redisReplicationSlaveOffset = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "slave_offset",
			Help:      "Slave id, IP address, port, offset",
		},
			[]string{
				"node_name",
				"node_address",
				"slave_id",
				"slave_addr",
			})
		redisReplicationSlaveLag = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "replication",
			Name:      "slave_lag",
			Help:      "Slave id, IP address, port, lag.",
		},
			[]string{
				"node_name",
				"node_address",
				"slave_id",
				"slave_addr",
			})
	)
	return &RedisReplicationCollector{
		redisReplicationRole,
		redisReplicationConnectedSlaves,
		redisReplicationMasterReplOffset,
		redisReplicationSecondReplOffset,
		redisReplicationReplBacklogActive,
		redisReplicationReplBacklogSize,
		redisReplicationReplBacklogFirstByteOffset,
		redisReplicationReplBacklogHistlen,
		redisReplicationMasterHostPort,
		redisReplicationMasterLinkStatus,
		redisReplicationMasterLastIOSecondsAgo,
		redisReplicationMasterSyncInProgress,
		redisReplicationSlaveReplOffset,
		redisReplicationSlavePriority,
		redisReplicationSlaveReadOnly,
		redisReplicationMasterSyncLeftBytes,
		redisReplicationMasterSyncLastIOSecondsAgo,
		redisReplicationMasterLinkDownSinceSeconds,
		redisReplicationMinSlavesGoodSlaves,
		redisReplicationSlaveState,
		redisReplicationSlaveOffset,
		redisReplicationSlaveLag,
	}
}

func (m *RedisReplicationCollector) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(m.Role)
	registry.MustRegister(m.ConnectedSlaves)
	//registry.MustRegister(info.RedisReplicationMasterRepl)
	registry.MustRegister(m.MasterReplOffset)
	registry.MustRegister(m.SecondReplOffset)
	registry.MustRegister(m.ReplBacklogActive)
	registry.MustRegister(m.ReplBacklogSize)
	registry.MustRegister(m.ReplBacklogFirstByteOffset)
	registry.MustRegister(m.ReplBacklogHistlen)
	registry.MustRegister(m.MasterHostPort)
	registry.MustRegister(m.MasterLinkStatus)
	registry.MustRegister(m.MasterLastIOSecondsAgo)
	registry.MustRegister(m.MasterSyncInProgress)
	registry.MustRegister(m.SlaveReplOffset)
	registry.MustRegister(m.SlavePriority)
	registry.MustRegister(m.SlaveReadOnly)
	registry.MustRegister(m.MasterSyncLeftBytes)
	registry.MustRegister(m.MasterSyncLastIOSecondsAgo)
	registry.MustRegister(m.MasterLinkDownSinceSeconds)
	registry.MustRegister(m.MinSlavesGoodSlaves)
	registry.MustRegister(m.SlaveState)
	registry.MustRegister(m.SlaveOffset)
	registry.MustRegister(m.SlaveLag)
}

func (m *RedisReplicationCollector) Unregister(registry *prometheus.Registry) bool {
	if !registry.Unregister(m.Role) {
		return false
	}
	if !registry.Unregister(m.ConnectedSlaves) {
		return false
	}
	//registry.Unregister(info.RedisReplicationMasterRepl)
	if !registry.Unregister(m.MasterReplOffset) {
		return false
	}
	if !registry.Unregister(m.SecondReplOffset) {
		return false
	}
	if !registry.Unregister(m.ReplBacklogActive) {
		return false
	}
	if !registry.Unregister(m.ReplBacklogSize) {
		return false
	}
	if !registry.Unregister(m.ReplBacklogFirstByteOffset) {
		return false
	}
	if !registry.Unregister(m.ReplBacklogHistlen) {
		return false
	}
	if !registry.Unregister(m.MasterHostPort) {
		return false
	}
	if !registry.Unregister(m.MasterLinkStatus) {
		return false
	}
	if !registry.Unregister(m.MasterLastIOSecondsAgo) {
		return false
	}
	if !registry.Unregister(m.MasterSyncInProgress) {
		return false
	}
	if !registry.Unregister(m.SlaveReplOffset) {
		return false
	}
	if !registry.Unregister(m.SlavePriority) {
		return false
	}
	if !registry.Unregister(m.SlaveReadOnly) {
		return false
	}
	if !registry.Unregister(m.MasterSyncLeftBytes) {
		return false
	}
	if !registry.Unregister(m.MasterSyncLastIOSecondsAgo) {
		return false
	}
	if !registry.Unregister(m.MasterLinkDownSinceSeconds) {
		return false
	}
	if !registry.Unregister(m.MinSlavesGoodSlaves) {
		return false
	}
	if !registry.Unregister(m.SlaveState) {
		return false
	}
	if !registry.Unregister(m.SlaveOffset) {
		return false
	}
	if !registry.Unregister(m.SlaveLag) {
		return false
	}
	return true
}

func (m *RedisReplicationCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	// role:master
	if role, ok := r["role"]; ok {
		switch role {
		case "master":
			m.Role.WithLabelValues(nodeName, nodeAddress).Set(float64(RedisReplicationRoleMaster))
		case "slave":
			m.Role.WithLabelValues(nodeName, nodeAddress).Set(float64(RedisReplicationRoleSlave))
		}
	}
	//connected_slaves:2
	if connectedSlavesStr, ok := r["connected_slaves"]; ok {
		if connectedSlaves, err := strconv.Atoi(connectedSlavesStr); err == nil {
			m.ConnectedSlaves.WithLabelValues(nodeName, nodeAddress).Set(float64(connectedSlaves))
		}
	}
	//master_repl_offset:2918675
	if masterReplOffsetStr, ok := r["master_repl_offset"]; ok {
		if masterReplOffset, err := strconv.Atoi(masterReplOffsetStr); err == nil {
			m.MasterReplOffset.WithLabelValues(nodeName, nodeAddress).Set(float64(masterReplOffset))
		}
	}
	//second_repl_offset:-1
	if secondReplOffsetStr, ok := r["second_repl_offset"]; ok {
		if secondReplOffset, err := strconv.Atoi(secondReplOffsetStr); err == nil {
			m.SecondReplOffset.WithLabelValues(nodeName, nodeAddress).Set(float64(secondReplOffset))
		}
	}
	//repl_backlog_active:1
	if replBacklogActiveStr, ok := r["repl_backlog_active"]; ok {
		if replBacklogActive, err := strconv.Atoi(replBacklogActiveStr); err == nil {
			m.ReplBacklogActive.WithLabelValues(nodeName, nodeAddress).Set(float64(replBacklogActive))
		}
	}
	//repl_backlog_size:1048576
	if replBacklogSizeStr, ok := r["repl_backlog_size"]; ok {
		if replBacklogSize, err := strconv.Atoi(replBacklogSizeStr); err == nil {
			m.ReplBacklogSize.WithLabelValues(nodeName, nodeAddress).Set(float64(replBacklogSize))
		}
	}
	//repl_backlog_first_byte_offset:1870100
	if replBacklogFirstByteOffsetStr, ok := r["repl_backlog_first_byte_offset"]; ok {
		if replBacklogFirstByteOffset, err := strconv.Atoi(replBacklogFirstByteOffsetStr); err == nil {
			m.ReplBacklogFirstByteOffset.WithLabelValues(nodeName, nodeAddress).Set(
				float64(replBacklogFirstByteOffset))
		}
	}
	//repl_backlog_histlen:1048576
	if replBacklogHistlenStr, ok := r["repl_backlog_histlen"]; ok {
		if replBacklogHistlen, err := strconv.Atoi(replBacklogHistlenStr); err == nil {
			m.ReplBacklogHistlen.WithLabelValues(nodeName, nodeAddress).Set(float64(replBacklogHistlen))
		}
	}
	// master_host:172.25.1.3
	// master_port:6379
	masterHost, hasMasterHost := r["masterHost"]
	masterPort, hasMasterPort := r["masterPort"]
	if hasMasterHost && hasMasterPort {
		m.MasterHostPort.WithLabelValues(nodeName, nodeAddress, masterHost, masterPort).Inc()
	}
	//master_link_status:up
	if masterLinkStatus, ok := r["master_link_status"]; ok {
		switch masterLinkStatus {
		case "up":
			m.MasterLinkStatus.WithLabelValues(nodeName, nodeAddress).Set(
				float64(RedisReplicationMasterLinkUp))
		case "down":
			m.MasterLinkStatus.WithLabelValues(nodeName, nodeAddress).Set(
				float64(RedisReplicationMasterLinkDown))
		}

	}
	//master_last_io_seconds_ago:1
	if masterLastIOSecondsAgoStr, ok := r["master_last_io_seconds_ago"]; ok {
		if masterLastIOSecondsAgo, err := strconv.Atoi(masterLastIOSecondsAgoStr); err == nil {
			m.MasterLastIOSecondsAgo.WithLabelValues(nodeName, nodeAddress).Set(
				float64(masterLastIOSecondsAgo))
		}
	}
	//master_sync_in_progress:0
	if masterSyncInProgressStr, ok := r["master_sync_in_progress"]; ok {
		if masterSyncInProgress, err := strconv.Atoi(masterSyncInProgressStr); err == nil {
			m.MasterSyncInProgress.WithLabelValues(nodeName, nodeAddress).Set(
				float64(masterSyncInProgress))
		}
	}
	//slave_repl_offset:5084610
	if slaveReplOffsetStr, ok := r["slave_repl_offset"]; ok {
		if slaveReplOffset, err := strconv.Atoi(slaveReplOffsetStr); err == nil {
			m.SlaveReplOffset.WithLabelValues(nodeName, nodeAddress).Set(float64(slaveReplOffset))
		}
	}
	//slave_priority:100
	if slavePriorityStr, ok := r["slave_priority"]; ok {
		if slavePriority, err := strconv.Atoi(slavePriorityStr); err == nil {
			m.SlavePriority.WithLabelValues(nodeName, nodeAddress).Set(float64(slavePriority))
		}
	}
	//slave_read_only:1
	if slaveReadOnlyStr, ok := r["slave_read_only"]; ok {
		if slaveReadOnly, err := strconv.Atoi(slaveReadOnlyStr); err == nil {
			m.SlaveReadOnly.WithLabelValues(nodeName, nodeAddress).Set(float64(slaveReadOnly))
		}
	}
	//slave0:ip=172.25.1.5,port=6379,state=online,offset=2918675,lag=0
	//slave1:ip=172.25.1.4,port=6379,state=online,offset=2918540,lag=1
	for key, value := range r {
		ok1 := strings.Contains(key, "slave")
		ok2 := strings.Contains(value, "state=")
		if ok1 && ok2 {
			valueMap := stringToMap(value)
			addr := fmt.Sprintf("%s:%s", valueMap["ip"], valueMap["port"])
			if state, ok := valueMap["state"]; ok {
				switch state {
				case "online":
					m.SlaveState.WithLabelValues(nodeName, nodeAddress, key, addr).Set(float64(
						RedisReplicationSlaveStateOnline))
				default:
					m.SlaveState.WithLabelValues(nodeName, nodeAddress, key, addr).Set(float64(
						RedisReplicationSlaveStateOffline))
				}
			}
			if offsetStr, ok := valueMap["offset"]; ok {
				if offset, err := strconv.Atoi(offsetStr); err == nil {
					m.SlaveOffset.WithLabelValues(nodeName, nodeAddress, key, addr).Set(float64(offset))
				}
			}
			if lagStr, ok := valueMap["lag"]; ok {
				if lag, err := strconv.Atoi(lagStr); err == nil {
					m.SlaveLag.WithLabelValues(nodeName, nodeAddress, key, addr).Set(float64(lag))
				}
			}

		}
	}

	return nil
}

func stringToMap(v string) map[string]string {
	m := make(map[string]string)
	values := strings.Split(v, ",")
	for _, item := range values {
		keyValue := strings.Split(item, "=")
		if len(keyValue) == 2 {
			m[keyValue[0]] = keyValue[1]
		}
	}
	return m
}
