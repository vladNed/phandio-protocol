package routes

import (
	"net/http"

	"github.com/mvx-mnr-atomic/signalling/internal/hub"
)

// serveWebSocket upgrades the HTTP connection to a WebSocket connection
// and registers the client with the hub.
func serveWebSocket(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Error upgrading to websocket: ", err)
		return
	}

	// Register new client
	client := hub.NewClient(h, conn)
	client.Register()

	// Start io streams
	go client.ReadStream()
	go client.WriteStream()
}