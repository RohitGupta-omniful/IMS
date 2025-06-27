package router

import (
	"github.com/RohitGupta-omniful/IMS/handler"
	"github.com/RohitGupta-omniful/IMS/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all service routes to the given server
func RegisterRoutes(s *gin.Engine) {
	// Hub CRUD
	s.Use(middleware.AuthMiddleware)
	s.POST("/hubs", handler.CreateHub)
	s.GET("/hubs", handler.ListHubs)
	//s.GET("/hubs/:id", handler.GetHub)
	s.GET("/validate/hub/:hub", handler.ValidateHubExists)
	s.PUT("/hubs/:id", handler.UpdateHub)
	s.DELETE("/hubs/:id", handler.DeleteHub)

	// SKU CRUD
	s.POST("/skus", handler.CreateSKU)
	s.GET("/skus", handler.ListSKUs)
	//s.GET("/skus/:id", handler.GetSKU)
	s.GET("/validate/sku/:sku", handler.ValidateSKUExists)
	s.PUT("/skus/:id", handler.UpdateSKU)
	s.DELETE("/skus/:id", handler.DeleteSKU)

	// Inventory
	s.POST("/inventory/update", handler.UpdateInventoryHandler)
	s.POST("/inventory/upsert", handler.UpsertInventory)
	s.GET("/inventory", handler.GetInventory)

}
