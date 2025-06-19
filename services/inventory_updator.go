package services

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ValidateHubExists(c *gin.Context) {
	hubIDStr := c.Param("id")
	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hub ID"})
		return
	}

	var hub models.Hub
	err = db.GetMasterDB(context.Background()).First(&hub, "id = ?", hubID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"exists": false})
		return
	} else if err != nil {
		log.Printf("[ValidateHubExists] DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true})
}

// --- HTTP Handler for validating SKU-on-Hub existence ---
func ValidateSKUOnHub(c *gin.Context) {
	skuIDStr := c.Query("sku_id")
	hubIDStr := c.Query("hub_id")

	skuID, err1 := uuid.Parse(skuIDStr)
	hubID, err2 := uuid.Parse(hubIDStr)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid SKU or Hub ID"})
		return
	}

	var inventory models.Inventory
	err := db.GetMasterDB(context.Background()).Where("product_id = ? AND hub_id = ?", skuID, hubID).First(&inventory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"exists": false})
		return
	} else if err != nil {
		log.Printf("[ValidateSKUOnHub] DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": true})
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
