package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/RohitGupta-omniful/IMS/config"
	commonsPostgres "github.com/omniful/go_commons/db/sql/postgres"
	commonsHTTP "github.com/omniful/go_commons/http"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var dbCluster *commonsPostgres.DbCluster

func main() {
	// Step 1: Load config
	ctx, err := config.LoadContext()
	if err != nil {
		log.Fatalf("Failed to load config context: %v", err)
	}
	cfg := config.AppConfig

	// Step 2: Init logger (optional, but good practice)
	// Step 2: Logger initialization skipped (logger package not available)
	// Step 3: Init DB
	initDatabase(ctx, cfg.Database)

	// Step 4: Init HTTP server using go_commons
	httpServer := commonsHTTP.InitializeServer(
		cfg.Server.Port,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		cfg.Server.IdleTimeout,
		false, // enable CORS
	)

	// Step 5: Register routes
	httpServer.POST("/users", createUserHandler)
	httpServer.GET("/users/first", getFirstUserHandler)

	// Step 6: Start the server
	log.Printf("Starting %s on %s...", cfg.Server.Name, cfg.Server.Port)
	if err := httpServer.StartServer(cfg.Server.Name); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initDatabase(ctx context.Context, cfg config.DBConfig) {
	dbConfig := commonsPostgres.DBConfig{
		Host:                   cfg.Host,
		Port:                   fmt.Sprintf("%d", cfg.Port),
		Username:               cfg.User,
		Password:               cfg.Password,
		Dbname:                 cfg.Name,
		DebugMode:              true,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		MaxOpenConnections:     10,
		MaxIdleConnections:     5,
		ConnMaxLifetime:        time.Minute * 10,
	}

	var slaveConfigs []commonsPostgres.DBConfig
	dbCluster = commonsPostgres.InitializeDBInstance(dbConfig, &slaveConfigs)

	db := dbCluster.GetMasterDB(ctx)
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}

func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbCluster.GetMasterDB(c.Request.Context())
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

func getFirstUserHandler(c *gin.Context) {
	var user User
	db := dbCluster.GetMasterDB(c.Request.Context())
	if err := db.First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}
