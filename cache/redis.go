package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RohitGupta-omniful/IMS/config"
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
			ClusterMode:                   config.GetRedisClusterMode(ctx),
			ServeReadsFromSlaves:          config.GetRedisServeReadsFromSlaves(ctx),
			ServeReadsFromMasterAndSlaves: config.GetRedisServeReadsFromMasterAndSlaves(ctx),
			PoolSize:                      uint(config.GetRedisPoolSize(ctx)),
			PoolFIFO:                      config.GetRedisPoolFIFO(ctx),
			MinIdleConn:                   uint(config.GetRedisMinIdleConn(ctx)),
			DB:                            uint(config.GetRedisDB(ctx)),
			Hosts:                         config.GetRedisHosts(ctx),
			DialTimeout:                   config.GetRedisDialTimeout(ctx),
			ReadTimeout:                   config.GetRedisReadTimeout(ctx),
			WriteTimeout:                  config.GetRedisWriteTimeout(ctx),
			IdleTimeout:                   config.GetRedisIdleTimeout(ctx),
		}

		client := redis.NewClient(cfg)
		if client == nil {
			log.Panic("Failed to initialize Redis client")
		}
		redisClient = client
		log.Println("Redis client initialized successfully")
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
	client := GetRedisClient()
	valStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}
	_, err := client.Set(ctx, key, valStr, ttl)
	return err
}

// Get retrieves a string value for a key
func Get(ctx context.Context, key string) (string, error) {
	client := GetRedisClient()
	return client.Get(ctx, key)
}

// SetJSON sets any struct as JSON in Redis
func SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	client := GetRedisClient()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = client.Set(ctx, key, string(data), ttl)
	return err
}

// GetJSON retrieves and unmarshals a JSON object into dest
func GetJSON(ctx context.Context, key string, dest interface{}) error {
	client := GetRedisClient()

	strVal, err := client.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(strVal), dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// Del deletes a key
func Del(ctx context.Context, keys ...string) error {
	client := GetRedisClient()
	_, err := client.Del(ctx, keys...)
	return err
}
