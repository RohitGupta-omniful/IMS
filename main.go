package main

import (
	"context"
	"fmt"
	"log"

	"github.com/RohitGupta-omniful/IMS/config"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http" // Using go_commons HTTP server
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *gorm.DB

func initDatabase(cfg config.DBConfig) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}

func main() {
	// Load config
	config.LoadConfig(".")
	cfg := config.AppConfig

	// Init DB
	initDatabase(cfg.Database)

	// Init HTTP server using go_commons
	server := http.InitializeServer(
		cfg.Server.Port,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		cfg.Server.IdleTimeout,
		false, // Set to true if you want to enable CORS or similar feature
	)

	// Routes
	server.POST("/users", createUserHandler)
	server.GET("/users/first", getFirstUserHandler)

	log.Printf("🚀 Starting %s on %s...", cfg.Server.Name, cfg.Server.Port)
	if err := server.StartServer(cfg.Server.Name); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := db.WithContext(context.Background()).Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, user)
}

func getFirstUserHandler(c *gin.Context) {
	var user User
	if err := db.WithContext(context.Background()).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}
