package db

import (
	"context"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"gorm.io/gorm"
)

var DBCluster *postgres.DbCluster

func InitDatabase(ctx context.Context) {
	// Load master DB config
	masterConfig := postgres.DBConfig{
		Host:                   config.GetString(ctx, "postgresql.master.host"),
		Port:                   config.GetString(ctx, "postgresql.master.port"),
		Username:               config.GetString(ctx, "postgresql.master.user"),
		Password:               config.GetString(ctx, "postgresql.master.password"),
		Dbname:                 config.GetString(ctx, "postgresql.database_name"),
		DebugMode:              config.GetBool(ctx, "postgresql.debugMode"),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		MaxOpenConnections:     config.GetInt(ctx, "postgresql.maxOpenConns"),
		MaxIdleConnections:     config.GetInt(ctx, "postgresql.maxIdleConns"),
		ConnMaxLifetime:        config.GetDuration(ctx, "postgresql.connMaxLifetime"),
	}

	// Load slave DB config
	slaveConfig := postgres.DBConfig{
		Host:     config.GetString(ctx, "postgresql.slave.host"),
		Port:     config.GetString(ctx, "postgresql.slave.port"),
		Username: config.GetString(ctx, "postgresql.slave.user"),
		Password: config.GetString(ctx, "postgresql.slave.password"),
		Dbname:   config.GetString(ctx, "postgresql.database_name"),
	}
	slaveConfigs := []postgres.DBConfig{slaveConfig}

	// Initialize DB cluster
	DBCluster = postgres.InitializeDBInstance(masterConfig, &slaveConfigs)
	log.Info(i18n.Translate(ctx, "Database connection established successfully"))
}

// GetMasterDB returns a master *gorm.DB instance
func GetMasterDB(ctx context.Context) *gorm.DB {
	return DBCluster.GetMasterDB(ctx)
}

// GetSlaveDB returns a slave *gorm.DB instance
func GetSlaveDB(ctx context.Context) *gorm.DB {
	return DBCluster.GetSlaveDB(ctx)
}
