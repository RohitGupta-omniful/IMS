package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationData struct {
	IsValid bool `json:"is_valid"`
}

type Response struct {
	Data ValidationData `json:"data"`
}

// AuthMiddleware checks for a valid Authorization header.
func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || token != "Bearer my-secret-token" {
		c.JSON(http.StatusBadRequest, Response{Data: ValidationData{IsValid: false}})
		c.Abort()
		return
	}
	c.Next()
}
