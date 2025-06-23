package main

import (
	"log"

	"github.com/RohitGupta-omniful/IMS/cache"
	"github.com/RohitGupta-omniful/IMS/config"
	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/db/migration"
	"github.com/RohitGupta-omniful/IMS/server"
)

func main() {
	// Load config context
	config.InitConfig()
	ctx, err := config.LoadContext()
	if err != nil {
		log.Fatalf("Failed to load config context: %v", err)
	}

	//Initialize database and Run database migrations
	db.InitDatabase(ctx)
	migration.RunMigrations(ctx)

	// Initialize server and routes
	app := server.Initialize(ctx)

	//redis cache client
	cache.InitRedisClient(ctx)

	// Start HTTP server
	log.Printf("Starting %s on %s", config.GetServerName(ctx), config.GetServerPort(ctx))
	if err := app.StartServer(config.GetServerName(ctx)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
