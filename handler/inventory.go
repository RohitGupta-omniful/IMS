package handler

import (
	"net/http"
	"time"

	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"gorm.io/gorm/clause"
)

func UpsertInventory(c *gin.Context) {
	var inv models.Inventory
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&inv); err != nil {
		log.Errorf("[UpsertInventory] %s: %v", i18n.Translate(ctx, "invalid_request"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "invalid_request")})
		return
	}

	dbConn := db.GetMasterDB(ctx)

	// Foreign key validation for hub_id
	var hub models.Hub
	if err := dbConn.First(&hub, "id = ?", inv.HubID).Error; err != nil {
		log.Errorf("[UpsertInventory] %s: %v", i18n.Translate(ctx, "invalid_hub_id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "invalid_hub_id")})
		return
	}

	// Foreign key validation for sku_id
	var sku models.SKU
	if err := dbConn.First(&sku, "id = ?", inv.ProductID).Error; err != nil {
		log.Errorf("[UpsertInventory] %s: %v", i18n.Translate(ctx, "invalid_sku_id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "invalid_sku_id")})
		return
	}

	inv.UpdatedAt = time.Now()

	tx := dbConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hub_id"}, {Name: "sku_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "updated_at"}),
	}).Create(&inv)

	if tx.Error != nil {
		log.Errorf("[UpsertInventory] %s: %v", i18n.Translate(ctx, "inventory_upsert_error"), tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(ctx, "inventory_upsert_error")})
		return
	}

	c.JSON(http.StatusOK, inv)
}

func GetInventory(c *gin.Context) {
	ctx := c.Request.Context()
	hubID := c.Query("hub_id")
	skuID := c.Query("sku_id")

	if hubID == "" || skuID == "" {
		log.Errorf("[GetInventory] %s", i18n.Translate(ctx, "missing_hub_or_sku"))
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "missing_hub_or_sku")})
		return
	}

	var inv models.Inventory
	err := db.GetMasterDB(ctx).First(&inv, "hub_id = ? AND sku_id = ?", hubID, skuID).Error

	if err != nil {
		log.Warnf("[GetInventory] Inventory not found for hub_id=%s, sku_id=%s", hubID, skuID)
		c.JSON(http.StatusOK, gin.H{
			"hub_id":   hubID,
			"sku_id":   skuID,
			"quantity": 0,
			"found":    false,
		})
		return
	}

	c.JSON(http.StatusOK, inv)
}
