package server

import (
	"context"

	"github.com/RohitGupta-omniful/IMS/config"
	"github.com/RohitGupta-omniful/IMS/router"
	"github.com/omniful/go_commons/http"
)

// Initialize sets up the routes and returns the initialized server
func Initialize(ctx context.Context) *http.Server {
	server := http.InitializeServer(
		config.GetServerPort(ctx),
		config.GetReadTimeout(ctx),
		config.GetWriteTimeout(ctx),
		config.GetIdleTimeout(ctx),
		false, // CORS disabled
	)
	router.RegisterRoutes(server.Engine)
	return server
}
