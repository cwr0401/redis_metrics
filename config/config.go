package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"gopkg.in/yaml.v2"
)

var (
	MaxDialTimeout time.Duration = 20 * time.Second
	MaxReadTimeout time.Duration = 20 * time.Second
)

type RedisAddress struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (a RedisAddress) address() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a RedisAddress) valid() error {
	if a.Port < 1 || a.Port > 65535 {
		return errors.New("the Redis port must between 1 and 65535")
	}
	return nil
}

type RedisConfig struct {
	Standalone []*RedisStandalone `yaml:"standalone,omitempty"`
	Sentinel   []*RedisSentinel   `yaml:"sentinel,omitempty"`
	Cluster    []*RedisCluster    `yaml:"cluster,omitempty"`
}

type RedisStandalone struct {
	Name         string `yaml:"name,omitempty"`
	RedisAddress `yaml:",inline"`
	Password     string        `yaml:"password,omitempty"`
	DialTimeout  time.Duration `yaml:"connect_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
}

func (r RedisStandalone) RedisOptions() *redis.Options {
	var options = redis.Options{}
	options.Addr = r.address()
	options.Network = "tcp"
	options.Password = r.Password
	options.DialTimeout = r.DialTimeout
	options.ReadTimeout = r.ReadTimeout
	return &options
}

type RedisSentinel struct {
	RedisStandalone   `yaml:",inline"`
	DiscoveryMaster   bool `yaml:"discovery_master"`
	DiscoverySlave    bool `yaml:"discovery_slave"`
	DiscoverySentinel bool `yaml:"discovery_sentinel"`
}

type RedisCluster struct {
	Name        string         `yaml:"name,omitempty"`
	Addrs       []RedisAddress `yaml:"addresses"`
	Password    string         `yaml:"password,omitempty"`
	DialTimeout time.Duration  `yaml:"connect_timeout"`
	ReadTimeout time.Duration  `yaml:"read_timeout"`
	Mode        string         `yaml:"mode"`
}

func (c RedisCluster) RedisClusterOptions() *redis.ClusterOptions {
	var clusterOptions = redis.ClusterOptions{}
	for _, addr := range c.Addrs {
		clusterOptions.Addrs = append(clusterOptions.Addrs, addr.address())
	}
	clusterOptions.Password = c.Password
	clusterOptions.DialTimeout = c.DialTimeout
	clusterOptions.ReadTimeout = c.ReadTimeout
	clusterOptions.ReadOnly = true
	return &clusterOptions
}

func ParseRedisConfig(f []byte) (*RedisConfig, error) {
	var redisConfig = RedisConfig{}
	err := yaml.Unmarshal(f, &redisConfig)
	if err != nil {
		return nil, err
	}
	for _, standalone := range redisConfig.Standalone {
		err = parseRedisStandalone(standalone)
		if err != nil {
			return nil, err
		}
	}
	for _, sentinel := range redisConfig.Sentinel {
		err = parseRedisSentinel(sentinel)
		if err != nil {
			return nil, err
		}
	}
	for _, cluster := range redisConfig.Cluster {
		err = parseRedisCluster(cluster)
		if err != nil {
			return nil, err
		}
	}
	return &redisConfig, nil
}

func parseRedisStandalone(s *RedisStandalone) error {
	if s.Port == 0 {
		s.Port = 6379
	}
	if s.Host == "" {
		s.Host = "localhost"
	}
	if s.DialTimeout.Seconds() > MaxDialTimeout.Seconds() {
		s.DialTimeout = MaxDialTimeout
	}
	if s.ReadTimeout.Seconds() > MaxReadTimeout.Seconds() {
		s.ReadTimeout = MaxReadTimeout
	}
	//if err := s.valid(); err != nil {
	//	return err
	//}
	return s.valid()
}

func parseRedisSentinel(s *RedisSentinel) error {
	if s.Port == 0 {
		s.Port = 26379
	}
	if s.Host == "" {
		s.Host = "localhost"
	}
	if s.DialTimeout.Seconds() > MaxDialTimeout.Seconds() {
		s.DialTimeout = MaxDialTimeout
	}
	if s.ReadTimeout.Seconds() > MaxReadTimeout.Seconds() {
		s.ReadTimeout = MaxReadTimeout
	}
	//if err := s.valid(); err != nil {
	//	return err
	//}
	return s.valid()
}

func parseRedisCluster(c *RedisCluster) error {
	for _, addr := range c.Addrs {
		if err := addr.valid(); err != nil {
			return err
		}
	}
	if c.DialTimeout.Seconds() > MaxDialTimeout.Seconds() {
		c.DialTimeout = MaxDialTimeout
	}
	if c.ReadTimeout.Seconds() > MaxReadTimeout.Seconds() {
		c.ReadTimeout = MaxReadTimeout
	}
	return nil
}
