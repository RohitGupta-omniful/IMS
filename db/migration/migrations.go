package migration

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/RohitGupta-omniful/IMS/config"
	"github.com/omniful/go_commons/db/sql/migration"
)

// RunMigrations initializes the migration and applies pending migrations
func RunMigrations(ctx context.Context) {
	// Extract DB config values
	username := config.GetMasterUser(ctx)
	rawPassword := config.GetMasterPassword(ctx)
	host := config.GetMasterHost(ctx)
	port := config.GetMasterPort(ctx)
	dbname := config.GetPostgresDBName(ctx)
	password := url.QueryEscape(rawPassword)

	// Construct correct DSN manually
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)
	log.Println("Migration DB URL:", dbURL)

	// Path to migration files
	migrationPath := "file://C:/Users/LENOVO/go/Omniful/Onboarding-project/InventoryManagementSystemMicroService/db/migration"

	// Initialize the migrator
	migrator, err := migration.InitializeMigrate(migrationPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	// Apply pending migrations
	if err := migrator.Up(); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
