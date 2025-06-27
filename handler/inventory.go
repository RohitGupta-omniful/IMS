package handler

import (
	"context"
	"time"

	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func UpsertInventory(c *gin.Context) {
	var inv models.Inventory

	if err := c.ShouldBindJSON(&inv); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	ctx := context.Background()
	dbConn := db.GetMasterDB(ctx)

	// Foreign key validation for hub_id
	var hub models.Hub
	if err := dbConn.First(&hub, "id = ?", inv.HubID).Error; err != nil {
		c.JSON(400, gin.H{"error": "invalid hub_id"})
		return
	}

	//Foreign key validation for sku_id
	var sku models.SKU
	if err := dbConn.First(&sku, "id = ?", inv.ProductID).Error; err != nil {
		c.JSON(400, gin.H{"error": "invalid sku_id"})
		return
	}

	inv.UpdatedAt = time.Now()

	tx := dbConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hub_id"}, {Name: "sku_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "updated_at"}),
	}).Create(&inv)

	if tx.Error != nil {
		c.JSON(500, gin.H{"error": "could not upsert inventory"})
		return
	}

	c.JSON(200, inv)
}

func GetInventory(c *gin.Context) {
	hubID := c.Query("hub_id")
	skuID := c.Query("sku_id")

	if hubID == "" || skuID == "" {
		c.JSON(400, gin.H{"error": "hub_id and sku_id are required"})
		return
	}

	ctx := context.Background()
	var inv models.Inventory
	err := db.GetMasterDB(ctx).First(&inv, "hub_id = ? AND sku_id = ?", hubID, skuID).Error

	// Default quantity to 0 if not found
	if err != nil {
		c.JSON(200, gin.H{
			"hub_id":   hubID,
			"sku_id":   skuID,
			"quantity": 0,
			"found":    false,
		})
		return
	}

	c.JSON(200, inv)
}
