package info

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type RedisCPUCollector struct {
	UsedCpuSys          *prometheus.GaugeVec
	UsedCpuUser         *prometheus.GaugeVec
	UsedCpuSysChildren  *prometheus.GaugeVec
	UsedCpuUserChildren *prometheus.GaugeVec
}

func NewRedisCPUCollector() *RedisCPUCollector {
	var (
		// CPU
		// used_cpu_sys:24.340000
		// used_cpu_user:5.900000
		// used_cpu_sys_children:0.000000
		// used_cpu_user_children:0.000000

		// used_cpu_sys:2352.870000
		redisCPUUsedCpuSys = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cpu",
			Name:      "used_cpu_sys",
			Help:      "System CPU consumed by the Redis server.",
		},
			[]string{"node_name", "node_address"})

		// used_cpu_user:624.490000
		redisCPUUsedCpuUser = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cpu",
			Name:      "used_cpu_user",
			Help:      "User CPU consumed by the Redis server.",
		},
			[]string{"node_name", "node_address"})

		// used_cpu_sys_children:0.000000
		redisCPUUsedCpuSysChildren = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cpu",
			Name:      "used_cpu_sys_children",
			Help:      "System CPU consumed by the background processes.",
		},
			[]string{"node_name", "node_address"})

		// used_cpu_user_children:0.000000
		redisCPUUsedCpuUserChildren = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "redis",
			Subsystem: "cpu",
			Name:      "used_cpu_user_children",
			Help:      "User CPU consumed by the background processes",
		},
			[]string{"node_name", "node_address"})
	)
	return &RedisCPUCollector{
		redisCPUUsedCpuSys,
		redisCPUUsedCpuUser,
		redisCPUUsedCpuSysChildren,
		redisCPUUsedCpuUserChildren,
	}
}

func (m *RedisCPUCollector) MustRegister(registry *prometheus.Registry) {
	registry.MustRegister(m.UsedCpuSys)
	registry.MustRegister(m.UsedCpuUser)
	registry.MustRegister(m.UsedCpuSysChildren)
	registry.MustRegister(m.UsedCpuUserChildren)
}

func (m *RedisCPUCollector) Unregister(registry *prometheus.Registry) bool {
	if !registry.Unregister(m.UsedCpuSys) {
		return false
	}
	if !registry.Unregister(m.UsedCpuUser) {
		return false
	}
	if !registry.Unregister(m.UsedCpuSysChildren) {
		return false
	}
	if !registry.Unregister(m.UsedCpuUserChildren) {
		return false
	}
	return true
}

func (m *RedisCPUCollector) Set(nodeName, nodeAddress string, r map[string]string) error {
	if usedCpuSysStr, ok := r["used_cpu_sys"]; ok {
		if usedCpuSys, err := strconv.ParseFloat(usedCpuSysStr, 64); err == nil {
			m.UsedCpuSys.WithLabelValues(nodeName, nodeAddress).Set(float64(usedCpuSys))
		}
	}
	if usedCpuUserStr, ok := r["used_cpu_user"]; ok {
		if usedCpuUser, err := strconv.ParseFloat(usedCpuUserStr, 64); err == nil {
			m.UsedCpuUser.WithLabelValues(nodeName, nodeAddress).Set(float64(usedCpuUser))
		}
	}
	if usedCpuSysChildrenStr, ok := r["used_cpu_sys_children"]; ok {
		if usedCpuSysChildren, err := strconv.ParseFloat(usedCpuSysChildrenStr, 64); err == nil {
			m.UsedCpuSysChildren.WithLabelValues(nodeName, nodeAddress).Set(float64(usedCpuSysChildren))
		}
	}
	if usedCpuUserChildrenStr, ok := r["used_cpu_user_children"]; ok {
		if usedCpuUserChildren, err := strconv.ParseFloat(usedCpuUserChildrenStr, 64); err == nil {
			m.UsedCpuUserChildren.WithLabelValues(nodeName, nodeAddress).Set(float64(usedCpuUserChildren))
		}
	}
	return nil
}
