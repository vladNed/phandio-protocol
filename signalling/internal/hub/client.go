package hub

import (
	"github.com/gorilla/websocket"
	"bytes"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client represents a single websocket connection in to the hub
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) Register() {
	c.hub.register <- c
}

func (c* Client) Unregister() {
	c.hub.unregister <- c
	c.conn.Close()
}

func (c *Client) ReadStream() {
	defer c.Unregister()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("Error reading message from websocket: ", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		logger.Info("Received message: ", string(message))
		c.hub.broadcast <- message
	}
}
