package db

import (
	"context"
	"log"

	"github.com/RohitGupta-omniful/IMS/config"
	"github.com/omniful/go_commons/db/sql/postgres"
	"gorm.io/gorm"
)

var DBCluster *postgres.DbCluster

// InitDatabase initializes the master and slave DB connections
func InitDatabase(ctx context.Context) {
	// Load master DB config
	masterConfig := postgres.DBConfig{
		Host:                   config.GetMasterHost(ctx),
		Port:                   config.GetMasterPort(ctx),
		Username:               config.GetMasterUser(ctx),
		Password:               config.GetMasterPassword(ctx),
		Dbname:                 config.GetPostgresDBName(ctx),
		DebugMode:              config.GetPostgresDebug(ctx),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		MaxOpenConnections:     config.GetPostgresMaxOpenConns(ctx),
		MaxIdleConnections:     config.GetPostgresMaxIdleConns(ctx),
		ConnMaxLifetime:        config.GetPostgresConnMaxLifetime(ctx),
	}

	// Load slave DB config
	slaveConfig := postgres.DBConfig{
		Host:     config.GetSlaveHost(ctx),
		Port:     config.GetSlavePort(ctx),
		Username: config.GetSlaveUser(ctx),
		Password: config.GetSlavePassword(ctx),
		Dbname:   config.GetPostgresDBName(ctx),
	}
	slaveConfigs := []postgres.DBConfig{slaveConfig}

	// Initialize cluster
	DBCluster = postgres.InitializeDBInstance(masterConfig, &slaveConfigs)

	log.Println("Database connection established successfully")
}

// GetMasterDB returns a master *gorm.DB instance
func GetMasterDB(ctx context.Context) *gorm.DB {
	return DBCluster.GetMasterDB(ctx)
}

// GetSlaveDB returns a slave *gorm.DB instance
func GetSlaveDB(ctx context.Context) *gorm.DB {
	return DBCluster.GetSlaveDB(ctx)
}
