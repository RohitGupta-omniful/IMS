package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"

	"InventoryManagementSystemMicroService/config"
)

func main() {
	config.LoadConfig(".")

	//db.Connect(config.AppConfig.Database)

	log.Printf("Starting %s at PORT %s\n", config.AppConfig.Server.Name, config.AppConfig.Server.Port)

	server := http.InitializeServer(
		config.AppConfig.Server.Port,
		config.AppConfig.Server.ReadTimeout,
		config.AppConfig.Server.WriteTimeout,
		config.AppConfig.Server.IdleTimeout,
		false,
	)

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	//ims := server.Engine.Group("/api/v1")
	//routes.RegisterHubRoutes(ims)

	if err := server.StartServer(config.AppConfig.Server.Name); err != nil {
		log.Fatal("Could not start server:", err)
	}
}
