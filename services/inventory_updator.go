package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RohitGupta-omniful/IMS/cache"
	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ValidationData struct {
	IsValid bool `json:"is_valid"`
}

type Response struct {
	Data ValidationData `json:"data"`
}

func ValidateHubExists(c *gin.Context) {
	hubIDStr := c.Param("hub")

	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Data: ValidationData{IsValid: false},
		})
		return
	}

	fmt.Println("reaches")
	ctx := context.Background()
	cacheKey := "hub:exists:" + hubID.String()

	// Try Redis first
	if cached, _ := cache.Get(ctx, cacheKey); cached != "" {
		exists := cached == "true"
		c.JSON(http.StatusOK, Response{
			Data: ValidationData{IsValid: exists},
		})
		//fmt.Println("found in redis")
		return
	}

	var hub models.Hub
	err = db.GetMasterDB(ctx).First(&hub, "id = ?", hubID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cache.Set(ctx, cacheKey, "false", 10*time.Minute)
		c.JSON(http.StatusOK, Response{
			Data: ValidationData{IsValid: false},
		})
		return
	} else if err != nil {
		log.Printf("[ValidateHubExists] DB error: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Data: ValidationData{IsValid: false},
		})
		return
	}
	//fmt.Println("found")
	cache.Set(ctx, cacheKey, "true", 10*time.Minute)
	c.JSON(http.StatusOK, Response{
		Data: ValidationData{IsValid: true},
	})
}

// --- HTTP Handler for validating SKU-on-Hub existence ---
func ValidateSKUExists(c *gin.Context) {
	skuIDStr := c.Param("sku")

	skuID, err := uuid.Parse(skuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Data: ValidationData{IsValid: false},
		})
		return
	}
	fmt.Println("reaches")
	ctx := context.Background()
	cacheKey := "sku:exists:" + skuID.String()

	if cached, _ := cache.Get(ctx, cacheKey); cached != "" {
		exists := cached == "true"
		c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: exists}})
		fmt.Println("found in redis")
		return
	}

	var skus models.SKU
	err = db.GetMasterDB(ctx).First(&skus, "id = ?", skuID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cache.Set(ctx, cacheKey, "false", 10*time.Minute)
		c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: false}})
		return
	} else if err != nil {
		log.Printf("[ValidateSKUExists] DB error: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Data: ValidationData{IsValid: false},
		})
		return
	}

	cache.Set(ctx, cacheKey, "true", 10*time.Minute)
	c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: true}})
}

var (
	ErrHubNotFound     = errors.New("hub not found")
	ErrSKUNotFound     = errors.New("SKU not found")
	ErrInsufficientQty = errors.New("not enough inventory")
)

// InventoryUpdateRequest represents the expected payload
type InventoryUpdateRequest struct {
	SKUID           uuid.UUID `json:"sku_id" binding:"required"`
	HubID           uuid.UUID `json:"hub_id" binding:"required"`
	QuantityChange  int       `json:"quantity_change" binding:"required"`
	TransactionType string    `json:"transaction_type" binding:"required"` // e.g., "add" or "remove"
}

// UpdateInventoryHandler handles inventory update HTTP requests
func UpdateInventoryHandler(c *gin.Context) {
	var req InventoryUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateInventoryHandler] Invalid payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := UpdateInventory(
		context.Background(),
		req.SKUID,
		req.HubID,
		req.QuantityChange,
		req.TransactionType,
	)

	if err != nil {
		log.Printf("[UpdateInventoryHandler] Service error: %v", err)
		switch err {
		case ErrHubNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "hub not found"})
		case ErrSKUNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		case ErrInsufficientQty:
			c.JSON(http.StatusBadRequest, gin.H{"error": "not enough inventory"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "inventory updated"})
}

// UpdateInventory performs the actual inventory update logic
func UpdateInventory(ctx context.Context, skuID, hubID uuid.UUID, quantityChange int, transactionType string) error {
	dbConn := db.GetMasterDB(ctx)

	// Validate hub exists
	var hub models.Hub
	if err := dbConn.First(&hub, "id = ?", hubID).Error; err != nil {
		return ErrHubNotFound
	}

	// Validate SKU exists
	var inventory models.Inventory
	if err := dbConn.Where("product_id = ? AND hub_id = ?", skuID, hubID).First(&inventory).Error; err != nil {
		return ErrSKUNotFound
	}

	// Update quantity logic
	newQty := inventory.Quantity + quantityChange
	if newQty < 0 {
		return ErrInsufficientQty
	}

	inventory.Quantity = newQty

	// Save updated inventory
	if err := dbConn.Save(&inventory).Error; err != nil {
		return err
	}

	return nil
}
