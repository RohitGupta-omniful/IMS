package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"strings"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/redis"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

// InitRedisClient initializes the Redis client once using go_commons/redis
func InitRedisClient(ctx context.Context) {
	once.Do(func() {
		cfg := &redis.Config{
			ClusterMode:                   config.GetBool(ctx, "redis.cluster_mode"),
			ServeReadsFromSlaves:          config.GetBool(ctx, "redis.serve_reads_from_slaves"),
			ServeReadsFromMasterAndSlaves: config.GetBool(ctx, "redis.serve_reads_from_master_and_slaves"),
			PoolSize:                      uint(config.GetInt(ctx, "redis.pool_size")),
			PoolFIFO:                      config.GetBool(ctx, "redis.pool_fifo"),
			MinIdleConn:                   uint(config.GetInt(ctx, "redis.min_idle_conn")),
			Hosts:                         strings.Split(config.GetString(ctx, "redis.hosts"), ","),
			DialTimeout:                   config.GetDuration(ctx, "redis.dial_timeout"),
			ReadTimeout:                   config.GetDuration(ctx, "redis.read_timeout"),
			WriteTimeout:                  config.GetDuration(ctx, "redis.write_timeout"),
			IdleTimeout:                   config.GetDuration(ctx, "redis.idle_timeout"),
		}

		client := redis.NewClient(cfg)
		if client == nil {
			log.Panic("Failed to initialize Redis client")
		}
		redisClient = client
		log.Info("Redis client initialized successfully")
	})
}

// GetRedisClient returns the initialized Redis client
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Panic("Redis client is not initialized. Call InitRedisClient first.")
	}
	return redisClient
}

// Set sets a key-value pair with expiration
func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	valStr, ok := value.(string)
	if !ok {
		log.Errorf("[Set] Value for key %s must be a string", key)
		return nil
	}
	_, err := GetRedisClient().Set(ctx, key, valStr, ttl)
	if err != nil {
		log.Errorf("[Set] Redis SET error for key %s: %v", key, err)
	}
	return err
}

// Get retrieves a string value for a key
func Get(ctx context.Context, key string) (string, error) {
	val, err := GetRedisClient().Get(ctx, key)
	if err != nil {
		log.Errorf("[Get] Redis GET error for key %s: %v", key, err)
	}
	return val, err
}

// SetJSON sets any struct as JSON in Redis
func SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Errorf("[SetJSON] JSON marshal error for key %s: %v", key, err)
		return err
	}
	_, err = GetRedisClient().Set(ctx, key, string(data), ttl)
	if err != nil {
		log.Errorf("[SetJSON] Redis SET error for key %s: %v", key, err)
	}
	return err
}

// GetJSON retrieves and unmarshals a JSON object into dest
func GetJSON(ctx context.Context, key string, dest interface{}) error {
	strVal, err := GetRedisClient().Get(ctx, key)
	if err != nil {
		log.Errorf("[GetJSON] Redis GET error for key %s: %v", key, err)
		return err
	}
	if err := json.Unmarshal([]byte(strVal), dest); err != nil {
		log.Errorf("[GetJSON] JSON unmarshal error for key %s: %v", key, err)
		return err
	}
	return nil
}

// Del deletes a key
func Del(ctx context.Context, keys ...string) error {
	_, err := GetRedisClient().Del(ctx, keys...)
	if err != nil {
		log.Errorf("[Del] Redis DEL error for keys %v: %v", keys, err)
	}
	return err
}
