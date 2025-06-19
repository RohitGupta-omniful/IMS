package handler

import (
	"context"
	"log"
	"time"

	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ----------------------- HUBS -----------------------

func CreateHub(c *gin.Context) {
	var hub models.Hub

	if err := c.BindJSON(&hub); err != nil {
		log.Printf("[CreateHub] Invalid JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
	}

	hub.ID = uuid.New()
	hub.CreatedAt = time.Now()
	hub.UpdatedAt = time.Now()

	if err := db.GetMasterDB(context.Background()).Create(&hub).Error; err != nil {
		log.Printf("[CreateHub] DB error: %v", err)
		c.JSON(500, gin.H{"error": "could not create hub"})
	}

	c.JSON(201, hub)
}

func GetHub(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": "invalid hub ID format"})
	}

	var hub models.Hub
	if err := db.GetMasterDB(context.Background()).First(&hub, "id = ?", id).Error; err != nil {
		log.Printf("[GetHub] Not found: %v", err)
		c.JSON(404, gin.H{"error": "hub not found"})
	}

	c.JSON(200, hub)
}

func ListHubs(c *gin.Context) {
	var hubs []models.Hub
	if err := db.GetMasterDB(context.Background()).Find(&hubs).Error; err != nil {
		log.Printf("[ListHubs] DB error: %v", err)
		c.JSON(500, gin.H{"error": "could not fetch hubs"})
	}
	c.JSON(200, hubs)
}

func UpdateHub(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": "invalid hub ID format"})
	}

	dbConn := db.GetMasterDB(context.Background())
	var existing models.Hub

	if err := dbConn.First(&existing, "id = ?", id).Error; err != nil {
		log.Printf("[UpdateHub] Not found: %v", err)
		c.JSON(404, gin.H{"error": "hub not found"})
	}

	var payload models.Hub
	if err := c.BindJSON(&payload); err != nil {
		log.Printf("[UpdateHub] Invalid JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
	}

	existing.Name = payload.Name
	existing.Location = payload.Location
	existing.UpdatedAt = time.Now()

	if err := dbConn.Save(&existing).Error; err != nil {
		log.Printf("[UpdateHub] Save error: %v", err)
		c.JSON(500, gin.H{"error": "could not update hub"})

	}

	c.JSON(200, existing)
}

func DeleteHub(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": "invalid hub ID format"})
	}

	if err := db.GetMasterDB(context.Background()).Delete(&models.Hub{}, "id = ?", id).Error; err != nil {
		log.Printf("[DeleteHub] Delete error: %v", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
	}

	c.JSON(200, map[string]string{"message": "hub deleted"})
}
