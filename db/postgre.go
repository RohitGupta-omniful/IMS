package db

import (
	"context"
	"log"

	"github.com/RohitGupta-omniful/IMS/config"
	"github.com/omniful/go_commons/db/sql/postgres"
)

var DBCluster *postgres.DbCluster

func InitDatabase(ctx context.Context) {
	dbConfig := postgres.DBConfig{
		Host:                   config.GetDBHost(ctx),
		Port:                   config.GetDBPort(ctx),
		Username:               config.GetDBUser(ctx),
		Password:               config.GetDBPassword(ctx),
		Dbname:                 config.GetDBName(ctx),
		DebugMode:              config.GetDBDebugMode(ctx),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		MaxOpenConnections:     config.GetDBMaxOpenConns(ctx),
		MaxIdleConnections:     config.GetDBMaxIdleConns(ctx),
		ConnMaxLifetime:        config.GetDBConnMaxLifetime(ctx),
	}

	var slaveConfigs []postgres.DBConfig
	slaveConfig := postgres.DBConfig{
		Host:     config.GetDBSlaveHost(ctx),
		Port:     config.GetDBSlavePort(ctx),
		Username: config.GetDBSlaveUser(ctx),
		Password: config.GetDBSlavePassword(ctx),
		Dbname:   config.GetDBName(ctx),
	}
	slaveConfigs = append(slaveConfigs, slaveConfig)

	DBCluster = postgres.InitializeDBInstance(dbConfig, &slaveConfigs)
	log.Println("✅ Database initialized successfully")
}

// GetMasterDB returns the master DB connection
func GetMasterDB(ctx context.Context) *postgres.PostgresDB {
	return DBCluster.GetMasterDB(ctx)
}

// GetSlaveDB returns the slave DB connection (optional)
func GetSlaveDB(ctx context.Context) *postgres.PostgresDB {
	return DBCluster.GetSlaveDB(ctx)
}
