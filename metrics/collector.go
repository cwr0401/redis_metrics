package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cwr0401/redis_metrics/config"
	"github.com/cwr0401/redis_metrics/metrics/info"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type RedisMetrics struct {
	Collector info.RedisCollector
	Duration  time.Duration
	Config    *config.RedisConfig
}

type RedisClient struct {
	Name   string
	Addr   string
	Client *redis.Client
}

func (r *RedisMetrics) Clients() []*RedisClient {
	var redisClients []*RedisClient
	for _, instance := range r.Config.RedisInstances {
		instanceOptions := instance.RedisOptions()
		name := instance.Name
		if instance.Name == "" {
			name = fmt.Sprintf("redis-server-%s", instanceOptions.Addr)
		}
		rdb := redis.NewClient(instanceOptions)
		rc := RedisClient{
			Name:   name,
			Addr:   instanceOptions.Addr,
			Client: rdb,
		}
		redisClients = append(redisClients, &rc)
	}

	return redisClients
}

func (r *RedisMetrics) Run(ctx context.Context, c chan<- struct{}) {
	wg := sync.WaitGroup{}
	clients := r.Clients()
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			for _, client := range clients {
				client.Client.Close()
			}
			c <- struct{}{}
			return
		default:
			for _, client := range clients {
				wg.Add(1)
				log.Infof("node=%s, addr=%s", client.Name, client.Addr)
				go func(client *RedisClient) {
					infoMap := redisInfoToMetrics(client.Name, client.Addr, client.Client)
					if err := r.Collector.Set(client.Name, client.Addr, infoMap); err != nil {
						log.WithFields(log.Fields{
							"node": client.Name,
							"addr": client.Addr,
						}).Errorf("Set Redis metrics error: %s", err)
					}
					wg.Done()
				}(client)
			}
			time.Sleep(r.Duration)
		}
	}
}
