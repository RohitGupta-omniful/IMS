package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/RohitGupta-omniful/IMS/cache"
	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSKU(c *gin.Context) {
	var sku models.SKU
	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sku.ID = uuid.New()
	sku.CreatedAt = time.Now()
	sku.UpdatedAt = time.Now()

	if err := db.GetMasterDB(context.Background()).Create(&sku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sku)
}

func ListSKUs(c *gin.Context) {
	ctx := context.Background()
	dbConn := db.GetMasterDB(ctx)

	var skus []models.SKU
	query := dbConn.Model(&models.SKU{})

	if tenantID := c.Query("tenant_id"); tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if sellerID := c.Query("seller_id"); sellerID != "" {
		query = query.Where("seller_id = ?", sellerID)
	}
	if skuCodes := c.QueryArray("sku_code"); len(skuCodes) > 0 {
		query = query.Where("sku_code IN ?", skuCodes)
	}

	if err := query.Find(&skus).Error; err != nil {
		c.JSON(500, gin.H{"error": "could not fetch SKUs"})
		return
	}
	c.JSON(200, skus)
}

func UpdateSKU(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	var sku models.SKU
	if err := db.GetMasterDB(ctx).First(&sku, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	}

	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sku.UpdatedAt = time.Now()

	if err := db.GetMasterDB(ctx).Save(&sku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Invalidate SKU cache
	_ = cache.Del(ctx, "sku:exists:"+id)

	c.JSON(http.StatusOK, sku)
}
func DeleteSKU(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	if err := db.GetMasterDB(ctx).Delete(&models.SKU{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Invalidate SKU cache
	_ = cache.Del(ctx, "sku:exists:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "SKU deleted"})
}
