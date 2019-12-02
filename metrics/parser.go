package metrics

import (
	"fmt"
	"strings"

	"github.com/cwr0401/redis_metrics/metrics/info"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"time"
)

func RedisInfoResultParser(infoResult string) map[string]string {
	infoResultLines := strings.Split(infoResult, "\r\n")
	infoResultMap := make(map[string]string)
	for _, line := range infoResultLines {
		if strings.ContainsRune(line, ':') {
			item := strings.SplitN(line, ":", 2)
			infoResultMap[item[0]] = item[1] // strings.Replace(item[1], "\r", "", -1)
		}
	}
	return infoResultMap
}

func redisInfoToMetrics(nodeName, nodeAddress string, rdb *redis.Client) error {
	err := rdb.Ping().Err()
	if err != nil {
		//log.Errorf("RedisStandalone %s(%s) is down.", name, options.Addr)
		//log.Error(err)
		info.RedisServerUp.WithLabelValues(nodeName, nodeAddress).Set(0)
		return err
	}
	info.RedisServerUp.WithLabelValues(nodeName, nodeAddress).Set(1)
	redisInfo, err := rdb.Info().Result()
	if err != nil {
		//log.Errorf("RedisStandalone %s(%s) execute info command fail.", name, options.Addr)
		//log.Error(err)
		return err
	}
	redisInfoMap := RedisInfoResultParser(redisInfo)

	// Server
	err = info.SetRedisServer(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Clients
	info.SetRedisClients(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Memory
	info.SetRedisMemory(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Persistence
	info.SetRedisPersistence(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Stats
	info.SetRedisStats(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Replication
	info.SetRedisReplication(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// CPU
	info.SetRedisCPU(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Cluster
	info.SetRedisCluster(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Keyspace
	info.SetRedisKeyspace(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	// Sentinel
	info.SetRedisSentinel(nodeName, nodeAddress, redisInfoMap)
	if err != nil {
		log.Errorf("")
	}
	return nil
}

func RedisServerMetrics(name string, options *redis.Options, duration time.Duration) {
	if name == "" {
		name = fmt.Sprintf("redis-standalone-%s", options.Addr)
	}
	rdb := redis.NewClient(options)
	defer rdb.Close()
	for {
		err := redisInfoToMetrics(name, options.Addr, rdb)
		if err != nil {
			log.Errorf("")
		}
		// interval
		time.Sleep(duration)
	}
}

func RedisSentinelMetrics(name string, options *redis.Options, duration time.Duration,
	discoveryMaster, discoverySlave, disvoerySentinel bool) {
	if name == "" {
		name = fmt.Sprintf("redis-sentinel-%s", options.Addr)
	}
	// self metrics
	go RedisServerMetrics(name, options, duration)
	// discovery
	if discoveryMaster || discoverySlave || disvoerySentinel {
		rs := redis.NewSentinelClient(options)
		defer rs.Close()
		for {
			err := rs.Ping().Err()
			if err != nil {
				log.Errorf("RedisSentinel %s(%s) is down.", name, options.Addr)
				log.Error(err)
				info.RedisSentinelUp.WithLabelValues(name, options.Addr).Set(0)
				time.Sleep(duration)
				continue
			}
			info.RedisSentinelUp.WithLabelValues(name, options.Addr).Set(0)

			mastersName, err := getSentinelMastersName(rs)
			if err != nil {
				time.Sleep(duration)
				continue
			}

			if discoveryMaster {
				for _, master := range mastersName {
					go redisSentinelMasterInfoToMetrics(name, master, rs)
				}
			}

			if discoverySlave {

			}

			if disvoerySentinel {

			}

			time.Sleep(duration)
		}
	}
}

func getSentinelMastersName(rs *redis.SentinelClient) ([]string, error) {
	masters, err := rs.Masters().Result()
	if err != nil {
		return nil, err
	}
	var mastersName []string
	for _, master := range masters {
		if master, ok := master.([]interface{}); ok {
			if nameKey, ok := master[0].(string); nameKey != "name" || !ok {
				continue
			}
			if nameValue, ok := master[1].(string); ok {
				mastersName = append(mastersName, nameValue)
			}
		}
	}
	return mastersName, nil
}

func redisSentinelMasterInfoToMetrics(nameNode, masterName string, rs *redis.SentinelClient) error {
	masterInfo, err := rs.Master(masterName).Result()
	if err != nil {
		return err
	}

	info.SetRedisSentinelMastersInfo(masterInfo)

	return nil
}
