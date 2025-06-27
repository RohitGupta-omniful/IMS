package handler

import (
	"context"
	"time"

	"github.com/RohitGupta-omniful/IMS/cache"
	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

func CreateHub(c *gin.Context) {
	var hub models.Hub
	ctx := context.Background()

	if err := c.BindJSON(&hub); err != nil {
		log.Errorf("[CreateHub] %s: %v", i18n.Translate(ctx, "invalid_request_body"), err)
		c.JSON(400, gin.H{"error": i18n.Translate(ctx, "invalid_request_body")})
		return
	}

	hub.ID = uuid.New()
	hub.CreatedAt = time.Now()
	hub.UpdatedAt = time.Now()

	if err := db.GetMasterDB(ctx).Create(&hub).Error; err != nil {
		log.Errorf("[CreateHub] %s: %v", i18n.Translate(ctx, "hub_create_failed"), err)
		c.JSON(500, gin.H{"error": i18n.Translate(ctx, "hub_create_failed")})
		return
	}

	c.JSON(201, hub)
}

func ListHubs(c *gin.Context) {
	ctx := context.Background()
	var hubs []models.Hub

	if err := db.GetMasterDB(ctx).Find(&hubs).Error; err != nil {
		log.Errorf("[ListHubs] %s: %v", i18n.Translate(ctx, "hubs_fetch_failed"), err)
		c.JSON(500, gin.H{"error": i18n.Translate(ctx, "hubs_fetch_failed")})
		return
	}

	c.JSON(200, hubs)
}

func UpdateHub(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": i18n.Translate(ctx, "invalid_hub_id_format")})
		return
	}

	dbConn := db.GetMasterDB(ctx)

	var existing models.Hub
	if err := dbConn.First(&existing, "id = ?", id).Error; err != nil {
		log.Errorf("[UpdateHub] %s: %v", i18n.Translate(ctx, "hub_not_found"), err)
		c.JSON(404, gin.H{"error": i18n.Translate(ctx, "hub_not_found")})
		return
	}

	var payload models.Hub
	if err := c.BindJSON(&payload); err != nil {
		log.Errorf("[UpdateHub] %s: %v", i18n.Translate(ctx, "invalid_request_body"), err)
		c.JSON(400, gin.H{"error": i18n.Translate(ctx, "invalid_request_body")})
		return
	}

	existing.Name = payload.Name
	existing.Location = payload.Location
	existing.UpdatedAt = time.Now()

	if err := dbConn.Save(&existing).Error; err != nil {
		log.Errorf("[UpdateHub] %s: %v", i18n.Translate(ctx, "hub_update_failed"), err)
		c.JSON(500, gin.H{"error": i18n.Translate(ctx, "hub_update_failed")})
		return
	}

	_ = cache.Del(ctx, "hub:exists:"+id)

	c.JSON(200, existing)
}

func DeleteHub(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(400, gin.H{"error": i18n.Translate(ctx, "invalid_hub_id_format")})
		return
	}

	if err := db.GetMasterDB(ctx).Delete(&models.Hub{}, "id = ?", id).Error; err != nil {
		log.Errorf("[DeleteHub] %s: %v", i18n.Translate(ctx, "hub_delete_failed"), err)
		c.JSON(400, gin.H{"error": i18n.Translate(ctx, "hub_delete_failed")})
		return
	}

	_ = cache.Del(ctx, "hub:exists:"+id)

	c.JSON(200, gin.H{"message": i18n.Translate(ctx, "hub_deleted")})
}
