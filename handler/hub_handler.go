package handler

import (
	"context"
	"time"

	"github.com/RohitGupta-omniful/IMS/cache"
	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/log"
)

// ----------------------- HUBS -----------------------

func CreateHub(c *gin.Context) {
	var hub models.Hub

	if err := c.BindJSON(&hub); err != nil {
		log.Errorf("[CreateHub] Invalid JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	hub.ID = uuid.New()
	hub.CreatedAt = time.Now()
	hub.UpdatedAt = time.Now()

	if err := db.GetMasterDB(context.Background()).Create(&hub).Error; err != nil {
		log.Errorf("[CreateHub] DB error: %v", err)
		c.JSON(500, gin.H{"error": "could not create hub"})
		return
	}

	c.JSON(201, hub)
}

func GetHub(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": "invalid hub ID format"})
		return
	}

	var hub models.Hub
	if err := db.GetMasterDB(context.Background()).First(&hub, "id = ?", id).Error; err != nil {
		log.Errorf("[GetHub] Not found: %v", err)
		c.JSON(404, gin.H{"error": "hub not found"})
		return
	}

	c.JSON(200, hub)
}

func ListHubs(c *gin.Context) {
	var hubs []models.Hub
	if err := db.GetMasterDB(context.Background()).Find(&hubs).Error; err != nil {
		log.Errorf("[ListHubs] DB error: %v", err)
		c.JSON(500, gin.H{"error": "could not fetch hubs"})
		return
	}
	c.JSON(200, hubs)
}

func UpdateHub(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": "invalid hub ID format"})
		return
	}

	ctx := context.Background()
	dbConn := db.GetMasterDB(ctx)

	var existing models.Hub
	if err := dbConn.First(&existing, "id = ?", id).Error; err != nil {
		log.Errorf("[UpdateHub] Not found: %v", err)
		c.JSON(404, gin.H{"error": "hub not found"})
		return
	}

	var payload models.Hub
	if err := c.BindJSON(&payload); err != nil {
		log.Errorf("[UpdateHub] Invalid JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	existing.Name = payload.Name
	existing.Location = payload.Location
	existing.UpdatedAt = time.Now()

	if err := dbConn.Save(&existing).Error; err != nil {
		log.Errorf("[UpdateHub] Save error: %v", err)
		c.JSON(500, gin.H{"error": "could not update hub"})
		return
	}

	// Invalidate Redis cache
	_ = cache.Del(ctx, "hub:exists:"+id)

	c.JSON(200, existing)
}

func DeleteHub(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": "invalid hub ID format"})
		return
	}

	ctx := context.Background()

	if err := db.GetMasterDB(ctx).Delete(&models.Hub{}, "id = ?", id).Error; err != nil {
		log.Errorf("[DeleteHub] Delete error: %v", err)
		c.JSON(400, gin.H{"error": "could not delete hub"})
		return
	}

	// Invalidate Redis cache
	_ = cache.Del(ctx, "hub:exists:"+id)

	c.JSON(200, map[string]string{"message": "hub deleted"})
}
