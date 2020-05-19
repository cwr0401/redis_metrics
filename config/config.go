package config

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"gopkg.in/yaml.v2"
	"time"
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
	RedisInstances RedisInstanceSlice `yaml:"redis,omitempty"`
}

func (r RedisConfig) Equal(n RedisConfig) bool {
	if len(r.RedisInstances) != len(n.RedisInstances) {
		return false
	}
	for i, r1 := range r.RedisInstances {
		r2 := n.RedisInstances[i]
		if r1 != r2 {
			return false
		}
	}
	return true
}

type RedisInstanceSlice []*RedisInstance

type RedisInstance struct {
	Name         string `yaml:"name,omitempty"`
	RedisAddress `yaml:",inline"`
	Password     string        `yaml:"password,omitempty"`
	DialTimeout  time.Duration `yaml:"connect_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
}

func (r *RedisInstance) RedisOptions() *redis.Options {
	var options = redis.Options{}
	options.Addr = r.address()
	options.Network = "tcp"
	options.Password = r.Password
	options.DialTimeout = r.DialTimeout
	options.ReadTimeout = r.ReadTimeout
	return &options
}

func ParseRedisConfig(f []byte) (*RedisConfig, error) {
	var redisConfig = RedisConfig{}
	err := yaml.Unmarshal(f, &redisConfig)
	if err != nil {
		return nil, err
	}
	for _, redisInstance := range redisConfig.RedisInstances {
		err = parseRedisInstance(redisInstance)
		if err != nil {
			return nil, err
		}
	}
	return &redisConfig, nil
}

func parseRedisInstance(s *RedisInstance) error {
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
