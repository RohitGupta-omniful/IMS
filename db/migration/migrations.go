package migration

import (
	"context"
	"net/url"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/migration"
	"github.com/omniful/go_commons/log"
)

func RunMigrations(ctx context.Context) {
	// Extract DB config values
	username := config.GetString(ctx, "postgresql.master.user")
	rawPassword := config.GetString(ctx, "postgresql.master.password")
	host := config.GetString(ctx, "postgresql.master.host")
	port := config.GetString(ctx, "postgresql.master.port")
	dbname := config.GetString(ctx, "postgresql.database_name")
	password := url.QueryEscape(rawPassword)

	// Construct DSN
	dbURL := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
	log.Infof("Migration DB URL: %s", dbURL)

	// Path to migration files (update to dynamic path or embed if needed)
	migrationPath := "file://C:/Users/LENOVO/go/Omniful/Onboarding-project/InventoryManagementSystemMicroService/db/migration"

	// Initialize migrator
	migrator, err := migration.InitializeMigrate(migrationPath, dbURL)
	if err != nil {
		log.Errorf("Failed to initialize migrator: %v", err)
		return
	}

	if err := migrator.Up(); err != nil {
		log.Errorf("Failed to apply migrations: %v", err)
		return
	}

	log.Info("Migrations applied successfully")
}
