package migrations

import (
	"fmt"
	"log"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/migration"
)

func RunMigration() {
	fmt.Println("Starting migration...")
	migrationPath := "file://C:/Users/Abhishek/Desktop/Omniful/OnboardingProject/IMS/migrations"

	ctx := mycontext.GetContext()
	myHost := config.GetString(ctx, "postgresql.master.host")
	myPort := config.GetString(ctx, "postgresql.master.port")
	myUsername := config.GetString(ctx, "postgresql.master.username")
	myPassword := config.GetString(ctx, "postgresql.master.password")
	myDbname := config.GetString(ctx, "postgresql.database")

	// 2. Build your database URL
	dbURL := migration.BuildSQLDBURL(
		myHost,     // host
		myPort,     // port
		myDbname,   // database name
		myUsername, // user
		myPassword, // password
	)

	// 3. Initialize the migrator
	migrator, err := migration.InitializeMigrate(migrationPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	// 4. Run the migration (up)
	if err := migrator.Up(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println(" Migration applied successfully!")
}
