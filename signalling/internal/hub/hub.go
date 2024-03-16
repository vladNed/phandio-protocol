package hub

import (
	"github.com/mvx-mnr-atomic/signalling/internal/logging"
)

var (
	logger = logging.GetLogger(nil)

	// Singleton instance of the hub. This is the only instance of the hub that
	// will be used by the application.
	HubInstance = NewHub()
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	logger.Info("Running hub instance.")
	for {
		select {
		case client := <-h.register:
			logger.Info("Registering client -> ", client.conn.RemoteAddr().String())
			h.clients[client] = true
		case client := <-h.unregister:
			logger.Info("Unregister client -> ", client.conn.RemoteAddr().String())
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
