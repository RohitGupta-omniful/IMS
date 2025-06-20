package config

import (
	"context"
	"log"
	"time"

	"github.com/omniful/go_commons/config"
)

const (
	// Server keys
	serverPortKey         = "server.port"
	serverNameKey         = "server.name"
	serverReadTimeoutKey  = "server.read_timeout"
	serverWriteTimeoutKey = "server.write_timeout"
	serverIdleTimeoutKey  = "server.idle_timeout"

	// PostgreSQL keys
	pgDebugKey           = "postgresql.debugMode"
	pgDBNameKey          = "postgresql.database_name"
	pgMaxIdleKey         = "postgresql.maxIdleConns"
	pgMaxOpenKey         = "postgresql.maxOpenConns"
	pgConnMaxLifetimeKey = "postgresql.connMaxLifetime"

	pgMasterHostKey     = "postgresql.master.host"
	pgMasterPortKey     = "postgresql.master.port"
	pgMasterUserKey     = "postgresql.master.user"
	pgMasterPasswordKey = "postgresql.master.password"

	pgSlaveHostKey     = "postgresql.slave.host"
	pgSlavePortKey     = "postgresql.slave.port"
	pgSlaveUserKey     = "postgresql.slave.user"
	pgSlavePasswordKey = "postgresql.slave.password"

	// AWS keys
	awsRegionKey        = "aws.region"
	awsPublicBucketKey  = "aws.public_bucket"
	awsPrivateBucketKey = "aws.private_bucket"

	//Redis keys
	redisHostsKey                      = "redis.hosts"
	redisClusterModeKey                = "redis.cluster_mode"
	redisServeReadsFromSlavesKey       = "redis.serve_reads_from_slaves"
	redisServeReadsFromMasterSlavesKey = "redis.serve_reads_from_master_and_slaves"
	redisPoolSizeKey                   = "redis.pool_size"
	redisPoolFIFOKey                   = "redis.pool_fifo"
	redisMinIdleConnKey                = "redis.min_idle_conn"
	redisDBKey                         = "redis.db"
	redisDialTimeoutKey                = "redis.dial_timeout"
	redisReadTimeoutKey                = "redis.read_timeout"
	redisWriteTimeoutKey               = "redis.write_timeout"
	redisIdleTimeoutKey                = "redis.idle_timeout"
)

// InitConfig initializes go_commons config with polling
func InitConfig() {
	err := config.Init(15 * time.Second)
	if err != nil {
		log.Panicf("Failed to initialize config: %v", err)
	}
}

// LoadContext returns a context with loaded config attached
func LoadContext() (context.Context, error) {
	ctx, err := config.TODOContext()
	if err != nil {
		log.Printf("Failed to load configuration context: %v", err)
		return nil, err
	}
	log.Println("Configuration context loaded successfully")
	return ctx, nil
}

// === Server ===

func GetServerName(ctx context.Context) string {
	return config.GetString(ctx, serverNameKey)
}

func GetServerPort(ctx context.Context) string {
	return config.GetString(ctx, serverPortKey)
}

func GetReadTimeout(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, serverReadTimeoutKey)
}

func GetWriteTimeout(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, serverWriteTimeoutKey)
}

func GetIdleTimeout(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, serverIdleTimeoutKey)
}

// === PostgreSQL ===

func GetPostgresDBName(ctx context.Context) string {
	return config.GetString(ctx, pgDBNameKey)
}

func GetPostgresDebug(ctx context.Context) bool {
	return config.GetBool(ctx, pgDebugKey)
}

func GetPostgresMaxIdleConns(ctx context.Context) int {
	return config.GetInt(ctx, pgMaxIdleKey)
}

func GetPostgresMaxOpenConns(ctx context.Context) int {
	return config.GetInt(ctx, pgMaxOpenKey)
}

func GetPostgresConnMaxLifetime(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, pgConnMaxLifetimeKey)
}

func GetMasterHost(ctx context.Context) string {
	return config.GetString(ctx, pgMasterHostKey)
}

func GetMasterPort(ctx context.Context) string {
	return config.GetString(ctx, pgMasterPortKey)
}

func GetMasterUser(ctx context.Context) string {
	return config.GetString(ctx, pgMasterUserKey)
}

func GetMasterPassword(ctx context.Context) string {
	return config.GetString(ctx, pgMasterPasswordKey)
}

func GetSlaveHost(ctx context.Context) string {
	return config.GetString(ctx, pgSlaveHostKey)
}

func GetSlavePort(ctx context.Context) string {
	return config.GetString(ctx, pgSlavePortKey)
}

func GetSlaveUser(ctx context.Context) string {
	return config.GetString(ctx, pgSlaveUserKey)
}

func GetSlavePassword(ctx context.Context) string {
	return config.GetString(ctx, pgSlavePasswordKey)
}

// === AWS ===

func GetAWSRegion(ctx context.Context) string {
	return config.GetString(ctx, awsRegionKey)
}

func GetPublicBucket(ctx context.Context) string {
	return config.GetString(ctx, awsPublicBucketKey)
}

func GetPrivateBucket(ctx context.Context) string {
	return config.GetString(ctx, awsPrivateBucketKey)
}

// === Redis ===

func GetRedisHosts(ctx context.Context) []string {
	return config.GetStringSlice(ctx, redisHostsKey)
}

func GetRedisClusterMode(ctx context.Context) bool {
	return config.GetBool(ctx, redisClusterModeKey)
}

func GetRedisServeReadsFromSlaves(ctx context.Context) bool {
	return config.GetBool(ctx, redisServeReadsFromSlavesKey)
}

func GetRedisServeReadsFromMasterAndSlaves(ctx context.Context) bool {
	return config.GetBool(ctx, redisServeReadsFromMasterSlavesKey)
}

func GetRedisPoolSize(ctx context.Context) int {
	return config.GetInt(ctx, redisPoolSizeKey)
}

func GetRedisPoolFIFO(ctx context.Context) bool {
	return config.GetBool(ctx, redisPoolFIFOKey)
}

func GetRedisMinIdleConn(ctx context.Context) int {
	return config.GetInt(ctx, redisMinIdleConnKey)
}

func GetRedisDB(ctx context.Context) int {
	return config.GetInt(ctx, redisDBKey)
}

func GetRedisDialTimeout(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, redisDialTimeoutKey)
}

func GetRedisReadTimeout(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, redisReadTimeoutKey)
}

func GetRedisWriteTimeout(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, redisWriteTimeoutKey)
}

func GetRedisIdleTimeout(ctx context.Context) time.Duration {
	return config.GetDuration(ctx, redisIdleTimeoutKey)
}
