package handler

import (
	"context"
	"net/http"
	"time"

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
	var skus []models.SKU
	if err := db.GetMasterDB(context.Background()).Find(&skus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, skus)
}

func GetSKU(c *gin.Context) {
	id := c.Param("id")
	var sku models.SKU
	if err := db.GetMasterDB(context.Background()).First(&sku, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	}
	c.JSON(http.StatusOK, sku)
}

func UpdateSKU(c *gin.Context) {
	id := c.Param("id")
	var sku models.SKU
	if err := db.GetMasterDB(context.Background()).First(&sku, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	}

	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sku.UpdatedAt = time.Now()
	if err := db.GetMasterDB(context.Background()).Save(&sku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sku)
}

func DeleteSKU(c *gin.Context) {
	id := c.Param("id")
	if err := db.GetMasterDB(context.Background()).Delete(&models.SKU{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "SKU deleted"})
}
