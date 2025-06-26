package server

import (
	"context"

	"github.com/RohitGupta-omniful/IMS/router"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
)

// Initialize sets up the routes and returns the initialized server
func Initialize(ctx context.Context) *http.Server {
	server := http.InitializeServer(
		config.GetString(ctx, "server.port"),
		config.GetDuration(ctx, "server.read_timeout"),
		config.GetDuration(ctx, "server.write_timeout"),
		config.GetDuration(ctx, "server.idle_timeout"),
		false,
	)
	router.RegisterRoutes(server.Engine)
	return server
}
