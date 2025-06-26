package main

import (
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"

	"github.com/RohitGupta-omniful/IMS/cache"
	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/db/migration"
	"github.com/RohitGupta-omniful/IMS/server"
)

func main() {
	// Initialize config
	if err := config.Init(15 * time.Second); err != nil {
		log.Errorf("Failed to initialize config: %v", err)
		return
	}

	// Get context with config
	ctx, err := config.TODOContext()
	if err != nil {
		log.Errorf("Failed to load config context: %v", err)
		return
	}

	// Initialize database and run migrations
	db.InitDatabase(ctx)
	migration.RunMigrations(ctx)

	// Initialize Redis cache
	cache.InitRedisClient(ctx)

	// Initialize HTTP server with routes
	app := server.Initialize(ctx)

	// Fetch server config values
	serverName := config.GetString(ctx, "server.name")
	serverPort := config.GetString(ctx, "server.port")

	log.Infof("Starting %s on %s", serverName, serverPort)

	// Start server
	if err := app.StartServer(serverPort); err != nil {
		log.Errorf("Server failed: %v", err)
	}
}
