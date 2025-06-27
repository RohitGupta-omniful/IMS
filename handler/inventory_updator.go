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
	"github.com/omniful/go_commons/i18n"
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
	ctx := c.Request.Context()
	hubIDStr := c.Param("hub")
	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		log.Errorf("[ValidateHubExists] %s: %v", i18n.Translate(ctx, "invalid_hub_id"), err)
		c.JSON(http.StatusBadRequest, Response{Data: ValidationData{IsValid: false}})
		return
	}

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
		log.Errorf("[ValidateHubExists] %s: %v", i18n.Translate(ctx, "db_error"), err)
		c.JSON(http.StatusInternalServerError, Response{Data: ValidationData{IsValid: false}})
		return
	}

	_ = cache.Set(ctx, cacheKey, "true", 10*time.Minute)
	c.JSON(http.StatusOK, Response{Data: ValidationData{IsValid: true}})
}

func ValidateSKUExists(c *gin.Context) {
	ctx := c.Request.Context()
	skuIDStr := c.Param("sku")
	skuID, err := uuid.Parse(skuIDStr)
	if err != nil {
		log.Errorf("[ValidateSKUExists] %s: %v", i18n.Translate(ctx, "invalid_sku_id"), err)
		c.JSON(http.StatusBadRequest, Response{Data: ValidationData{IsValid: false}})
		return
	}

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
		log.Errorf("[ValidateSKUExists] %s: %v", i18n.Translate(ctx, "db_error"), err)
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
	ctx := c.Request.Context()

	var req InventoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("[UpdateInventoryHandler] %s: %v", i18n.Translate(ctx, "invalid_request"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "invalid_request")})
		return
	}

	skuUUID, err := uuid.Parse(req.SKUID)
	if err != nil {
		log.Errorf("[UpdateInventoryHandler] %s: %v", i18n.Translate(ctx, "invalid_sku_id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "invalid_sku_id")})
		return
	}

	hubUUID, err := uuid.Parse(req.HubID)
	if err != nil {
		log.Errorf("[UpdateInventoryHandler] %s: %v", i18n.Translate(ctx, "invalid_hub_id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "invalid_hub_id")})
		return
	}

	err = UpdateInventory(ctx, skuUUID, hubUUID, req.QuantityChange, req.TransactionType)
	if err != nil {
		log.Errorf("[UpdateInventoryHandler] %s: %v", i18n.Translate(ctx, "inventory_update_error"), err)
		switch err {
		case ErrHubNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(ctx, "hub_not_found")})
		case ErrSKUNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(ctx, "sku_not_found")})
		case ErrInsufficientQty:
			c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(ctx, "insufficient_inventory")})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(ctx, "internal_server_error")})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n.Translate(ctx, "inventory_updated")})
}

func UpdateInventory(ctx context.Context, skuID, hubID uuid.UUID, quantityChange int, transactionType string) error {
	dbConn := db.GetMasterDB(ctx)

	return dbConn.Transaction(func(tx *gorm.DB) error {
		var inventory models.Inventory

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
			// no-op
		case "remove":
			quantityChange = -intAbs(quantityChange)
		default:
			return errors.New(i18n.Translate(ctx, "invalid_transaction_type"))
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
