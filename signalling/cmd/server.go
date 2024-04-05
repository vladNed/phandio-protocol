package main

import (
	"net/http"

	"github.com/mvx-mnr-atomic/signalling/internal/logging"
	"github.com/mvx-mnr-atomic/signalling/internal/settings"
	"github.com/mvx-mnr-atomic/signalling/internal/routes"
	"github.com/mvx-mnr-atomic/signalling/internal/hub"
)

var(
	config = settings.GetSettings()
	logger = logging.GetLogger(nil)
)

func main() {
	// Initialize hub
	go hub.HubInstance.Run()

	// Register routes
	mux := http.NewServeMux()
	routes.HttpRouter.Register(mux)
	routes.WSRouter.Register(mux)

	// Start server
	logger.Info("Running server on: ", config.GetAddress())
	err := http.ListenAndServeTLS(config.GetAddress(), config.CertFile, config.KeyFile, mux)
	if err != nil {
		logger.Error("Error starting server: ", err)
	}
}
