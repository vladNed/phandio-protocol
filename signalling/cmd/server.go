package main

import (
	"net/http"

	"github.com/mvx-mnr-atomic/signalling/internal/logging"
	"github.com/mvx-mnr-atomic/signalling/internal/settings"
	"github.com/mvx-mnr-atomic/signalling/internal/routes"
	"github.com/mvx-mnr-atomic/signalling/internal/hub"
)

var config = settings.GetSettings()
var logger = logging.GetLogger(nil)

func main() {
	// Initialize hub
	go hub.HubInstance.Run()


	// Register routes
	routes.HttpRouter.Register()
	routes.WSRouter.Register()

	// Start server
	logger.Info("Running server on: ", config.GetAddress())
	http.ListenAndServe(config.GetAddress(), nil)
}
