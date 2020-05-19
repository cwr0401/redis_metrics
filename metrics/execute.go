package metrics

import (
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

func redisInfoToMetrics(nodeName, nodeAddress string, rdb *redis.Client) map[string]string {
	log.WithFields(log.Fields{
		"node": nodeName,
		"addr": nodeAddress,
	}).Info("Get Redis Server Info Metrics")

	err := rdb.Ping().Err()
	if err != nil {
		log.WithFields(log.Fields{
			"node": nodeName,
			"addr": nodeAddress,
		}).Error("Redis Server is down.")
		log.WithFields(log.Fields{
			"node": nodeName,
			"addr": nodeAddress,
		}).Error(err)
		return map[string]string{
			"down": "",
		}
	}

	redisInfo, err := rdb.Info("all").Result()

	if err != nil {
		log.WithFields(log.Fields{
			"node": nodeName,
			"addr": nodeAddress,
		}).Errorf("Execute command 'info all' failed: %s", err)
		return map[string]string{
			"down": "",
		}
	}
	redisInfoMap := RedisInfoResultParser(redisInfo)
	return redisInfoMap
}
