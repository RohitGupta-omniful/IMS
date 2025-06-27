package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/RohitGupta-omniful/IMS/cache"
	"github.com/RohitGupta-omniful/IMS/db"
	"github.com/RohitGupta-omniful/IMS/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		c.JSON(http.StatusBadRequest, Response{Data: ValidationData{IsValid: false}})
		return
	}

	ctx := context.Background()
	cacheKey := "hub:exists:" + hubID.String()

	if cached, _ := cache.Get(ctx, cacheKey); cached != "" {
		c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: cached == "true"}})
		return
	}

	var hub models.Hub
	err = db.GetMasterDB(ctx).First(&hub, "id = ?", hubID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = cache.Set(ctx, cacheKey, "false", 10*time.Minute)
		c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: false}})
		return
	} else if err != nil {
		log.Errorf("[ValidateHubExists] DB error: %v", err)
		c.JSON(http.StatusInternalServerError, Response{Data: ValidationData{IsValid: false}})
		return
	}

	_ = cache.Set(ctx, cacheKey, "true", 10*time.Minute)
	c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: true}})
}

func ValidateSKUExists(c *gin.Context) {
	skuIDStr := c.Param("sku")
	skuID, err := uuid.Parse(skuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Data: ValidationData{IsValid: false}})
		return
	}

	ctx := context.Background()
	cacheKey := "sku:exists:" + skuID.String()

	if cached, _ := cache.Get(ctx, cacheKey); cached != "" {
		c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: cached == "true"}})
		return
	}

	var sku models.SKU
	err = db.GetMasterDB(ctx).First(&sku, "id = ?", skuID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = cache.Set(ctx, cacheKey, "false", 10*time.Minute)
		c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: false}})
		return
	} else if err != nil {
		log.Errorf("[ValidateSKUExists] DB error: %v", err)
		c.JSON(http.StatusInternalServerError, Response{Data: ValidationData{IsValid: false}})
		return
	}

	_ = cache.Set(ctx, cacheKey, "true", 10*time.Minute)
	c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: true}})
}

var (
	ErrHubNotFound     = errors.New("hub not found")
	ErrSKUNotFound     = errors.New("SKU not found")
	ErrInsufficientQty = errors.New("not enough inventory")
)

type InventoryUpdateRequest struct {
	SKUID           string `json:"sku_id"`
	HubID           string `json:"hub_id"`
	QuantityChange  int    `json:"quantity_change"`
	TransactionType string `json:"transaction_type"`
}

func UpdateInventoryHandler(c *gin.Context) {
	var req InventoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("[UpdateInventoryHandler] Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": Translate("invalid_request")})
		return
	}

	skuUUID, err := uuid.Parse(req.SKUID)
	if err != nil {
		log.Errorf("[UpdateInventoryHandler] Invalid sku_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": Translate("invalid_sku_id")})
		return
	}

	hubUUID, err := uuid.Parse(req.HubID)
	if err != nil {
		log.Errorf("[UpdateInventoryHandler] Invalid hub_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": Translate("invalid_hub_id")})
		return
	}

	err = UpdateInventory(context.Background(), skuUUID, hubUUID, req.QuantityChange, req.TransactionType)
	if err != nil {
		log.Errorf("[UpdateInventoryHandler] Service error: %v", err)
		switch err {
		case ErrHubNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": Translate("hub_not_found")})
		case ErrSKUNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": Translate("sku_not_found")})
		case ErrInsufficientQty:
			c.JSON(http.StatusBadRequest, gin.H{"error": Translate("insufficient_inventory")})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": Translate("internal_server_error")})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": Translate("inventory_updated")})
}

func UpdateInventory(ctx context.Context, skuID, hubID uuid.UUID, quantityChange int, transactionType string) error {
	dbConn := db.GetMasterDB(ctx)

	return dbConn.Transaction(func(tx *gorm.DB) error {
		var inventory models.Inventory

		// Lock the row to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("sku_id = ? AND hub_id = ?", skuID, hubID).
			First(&inventory).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSKUNotFound
			}
			return err
		}

		switch transactionType {
		case "add":
			// quantityChange remains positive
		case "remove":
			quantityChange = -intAbs(quantityChange)
		default:
			return errors.New(Translate("invalid_transaction_type"))
		}

		newQty := inventory.Quantity + quantityChange
		if newQty < 0 {
			return ErrInsufficientQty
		}

		inventory.Quantity = newQty
		return tx.Save(&inventory).Error
	})
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// i18n stub
var locale = map[string]string{
	"invalid_request":          "Invalid request body",
	"invalid_sku_id":           "Invalid SKU ID",
	"invalid_hub_id":           "Invalid Hub ID",
	"hub_not_found":            "Hub not found",
	"sku_not_found":            "SKU not found",
	"insufficient_inventory":   "Not enough inventory",
	"internal_server_error":    "Internal server error",
	"inventory_updated":        "Inventory updated successfully",
	"invalid_transaction_type": "Invalid transaction type",
}

func Translate(key string) string {
	if val, ok := locale[key]; ok {
		return val
	}
	return key
}
