package hub

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mvx-mnr-atomic/signalling/internal/cache"
)

var (
	writeWait = 1 * time.Second
	newline   = []byte{'\n'}
	space     = []byte{' '}
)

// Client represents a single websocket connection in to the hub
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send     chan []byte
	channels []string
	state    State
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:   hub,
		conn:  conn,
		send:  make(chan []byte, 256),
		state: New,
	}
}

func (c *Client) Register() {
	c.hub.register <- c
}

func (c *Client) Unregister() {
	c.hub.unregister <- c
	c.conn.Close()
}

func (c *Client) handleChannelSubscribe(message []byte) {
	channelSubscribeRequest, err := parseChannelRequest(message)
	if err != nil {
		logger.Error("Error parsing channel request: ", err)
		return
	}
	c.channels = channelSubscribeRequest.Channels
	logger.Info("Client subscribed to channels: ", channelSubscribeRequest.Channels)
	response := MessageResponse{Status: 200, Details: "Subscribed to channels"}
	responsePayload, err := parseMessageResponse(response)
	if err != nil {
		logger.Error("Error parsing channel response: ", err)
		return
	}
	c.state = Registered
	c.send <- responsePayload
}

func (c *Client) handleMessage(message []byte) {
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	messageRequest, err := parseMessageRequest(message)
	if err != nil {
		logger.Error("Error parsing message request: ", err)
		return
	}

	switch messageRequest := messageRequest.(type) {
	case *CreateOfferRequest:
		c.state = OfferCreated
		cache.MemcacheInstance.Set(messageRequest.OfferID, c)
		response := MessageResponse{Status: 200}
		responsePayload, err := parseMessageResponse(response)
		if err != nil {
			logger.Error("Error parsing message response: ", err)
			break
		}
		c.send <- responsePayload
		c.hub.broadcast <- message
		logger.Info("New offer created: ", messageRequest.OfferID)
	case *AnswerOfferRequest:
		if c.state != Registered {
			response := MessageResponse{Status: 404, Details: "Invalid state for answer request"}
			responsePayload, err := parseMessageResponse(response)
			if err != nil {
				logger.Error("Error parsing message response: ", err)
				break
			}
			c.send <- responsePayload
			return
		}
		if !cache.MemcacheInstance.Contains(messageRequest.OfferID) {
			response := MessageResponse{Status: 404, Details: "Offer not found"}
			responsePayload, err := parseMessageResponse(response)
			if err != nil {
				logger.Error("Error parsing message response: ", err)
				break
			}
			c.send <- responsePayload
			return
		}
		client, ok := cache.MemcacheInstance.Get(messageRequest.OfferID)
		if !ok {
			logger.Error("Error casting client from cache")
			return
		}
		cl := client.(*Client)
		cl.send <- message
		c.state = OfferAccepted

		response := MessageResponse{Status: 200}
		responsePayload, err := parseMessageResponse(response)
		if err != nil {
			logger.Error("Error parsing message response: ", err)
			break
		}
		c.send <- responsePayload
	default:
		logger.Error("Invalid message type")
	}
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

		if c.state == New {
			c.handleChannelSubscribe(message)
			continue
		}

		c.handleMessage(message)
	}
}

func (c *Client) WriteStream() {
	defer c.Unregister()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				continue
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		default:
			continue
		}
	}
}
