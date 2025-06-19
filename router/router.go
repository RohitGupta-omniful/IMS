package router

import (
	"github.com/RohitGupta-omniful/IMS/handler"
	"github.com/RohitGupta-omniful/IMS/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all service routes to the given server
func RegisterRoutes(s *gin.Engine) {
	// Hub CRUD
	s.POST("/hubs", handler.CreateHub)
	s.GET("/hubs", handler.ListHubs)
	s.GET("/hubs/:id", handler.GetHub)
	s.PUT("/hubs/:id", handler.UpdateHub)
	s.DELETE("/hubs/:id", handler.DeleteHub)

	// SKU CRUD
	s.POST("/skus", handler.CreateSKU)
	s.GET("/skus", handler.ListSKUs)
	s.GET("/skus/:id", handler.GetSKU)
	s.PUT("/skus/:id", handler.UpdateSKU)
	s.DELETE("/skus/:id", handler.DeleteSKU)

	// Inventory update
	s.POST("/inventory/update", services.UpdateInventoryHandler)
	s.GET("/validate/hub", services.ValidateHubExists)
	s.GET("/validate/sku_on_hub", services.ValidateSKUOnHub)
}
