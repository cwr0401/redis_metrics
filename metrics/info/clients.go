package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var (
	// Clients
	// connected_clients
	RedisClientsConnectedClients = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "clients",
		Name:      "connected_clients",
		Help:      "Number of client connections (excluding connections from replicas).",
	},
		[]string{"node_name", "node_address"})

	// client_recent_max_input_buffer
	// client_biggest_input_buf
	RedisClientsClientRecentMaxInputBuffer = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "clients",
		Name:      "client_recent_max_input_buffer",
		Help:      "biggest input buffer among current client connections",
	},
		[]string{"node_name", "node_address"})

	// client_recent_max_output_buffer
	// client_longest_output_list
	RedisClientsClientRecentMaxOutputBuffer = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "clients",
		Name:      "client_recent_max_output_buffer",
		Help:      "longest output list among current client connections.",
	},
		[]string{"node_name", "node_address"})

	// client_recent_max_input_buffer
	//RedisClientsClientRecentMaxInputBuffer = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	//	Namespace: "redis",
	//	Subsystem: "clients",
	//	Name:      "client_recent_max_input_buffer",
	//	Help:      "",
	//},
	//	[]string{"node_name", "node_address"})

	//client_recent_max_output_buffer
	//RedisClientsClientRecentMaxOutputBuffer = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	//	Namespace: "redis",
	//	Subsystem: "clients",
	//	Name:      "client_recent_max_output_buffer",
	//	Help:      "",
	//},
	//	[]string{"node_name", "node_address"})

	// blocked_clients
	RedisClientsBlockedClients = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "redis",
		Subsystem: "clients",
		Name:      "blocked_clients",
		Help:      "Number of clients pending on a blocking call (BLPOP, BRPOP, BRPOPLPUSH)",
	},
		[]string{"node_name", "node_address"})
)

func SetRedisClients(nodeName, nodeAddress string, r map[string]string) error {
	if connectedClientsStr, ok := r["connected_clients"]; ok {
		if connectedClients, err := strconv.Atoi(connectedClientsStr); err == nil {
			RedisClientsConnectedClients.WithLabelValues(nodeName, nodeAddress).Set(float64(connectedClients))
		}
	}
	if clientRecentMaxInputBufferStr, ok := r["client_recent_max_input_buffer"]; ok {
		if clientRecentMaxInputBuffer, err := strconv.Atoi(clientRecentMaxInputBufferStr); err == nil {
			RedisClientsClientRecentMaxInputBuffer.WithLabelValues(
				nodeName, nodeAddress,
			).Set(float64(clientRecentMaxInputBuffer))
		}
	} else {
		if ClientBiggestInputBufStr, ok := r["client_biggest_input_buf"]; ok {
			if ClientBiggestInputBuf, err := strconv.Atoi(ClientBiggestInputBufStr); err == nil {
				RedisClientsClientRecentMaxInputBuffer.WithLabelValues(
					nodeName, nodeAddress,
				).Set(float64(ClientBiggestInputBuf))
			}
		}
	}

	if ClientRecentMaxOutputBufferStr, ok := r["client_recent_max_output_buffer"]; ok {
		if ClientRecentMaxOutputBuffer, err := strconv.Atoi(ClientRecentMaxOutputBufferStr); err == nil {
			RedisClientsClientRecentMaxOutputBuffer.WithLabelValues(
				nodeName, nodeAddress,
			).Set(float64(ClientRecentMaxOutputBuffer))
		}
	} else {
		if clientLongestOutputListStr, ok := r["client_longest_output_list"]; ok {
			if clientLongestOutputList, err := strconv.Atoi(clientLongestOutputListStr); err == nil {
				RedisClientsClientRecentMaxOutputBuffer.WithLabelValues(
					nodeName, nodeAddress,
				).Set(float64(clientLongestOutputList))
			}
		}
	}

	if BlockedClientsStr, ok := r["blocked_clients"]; ok {
		if BlockedClients, err := strconv.Atoi(BlockedClientsStr); err == nil {
			RedisClientsBlockedClients.WithLabelValues(nodeName, nodeAddress).Set(float64(BlockedClients))
		}
	}
	return nil
}
