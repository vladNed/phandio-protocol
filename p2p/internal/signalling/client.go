package signalling

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"nhooyr.io/websocket"

	"github.com/mvx-mnr-atomic/p2p/internal/peer"
)

var (
	signallingServer = "ws://localhost:8080/ws/v1/" // TODO: Move to config
)

type WSClient struct {
	conn *websocket.Conn
}

func NewWSClient() (*WSClient, error) {
	conn, _, err := websocket.Dial(context.Background(), signallingServer, nil)
	if err != nil {
		log.Println("Could not connect to websocket. err:", err)
		return nil, err
	}

	return &WSClient{conn: conn}, nil
}

func (ws *WSClient) RegisterChannels() bool {
	channelsReq := RegisterChannelRequest{Channels: []string{"offers"}}
	payload, err := json.Marshal(channelsReq)
	if err != nil {
		return false
	}

	err = ws.conn.Write(context.Background(), websocket.MessageText, payload)
	if err != nil {
		return false
	}

	msgType, payload, err := ws.conn.Read(context.Background())
	if err != nil {
		return false
	}

	if msgType != websocket.MessageText {
		return false
	}

	resp, err := parseResponse(payload)
	if err != nil {
		return false
	}

	if resp.Status != WS_OK_STATUS {
		return false
	}

	return true
}

func (ws *WSClient) Listener(
	pageLogger func(string),
	peer *peer.Peer,
	marketplace *Marketplace,
) {
	defer ws.Close()
	for {
		msgType, payload, err := ws.conn.Read(context.Background())
		if err != nil {
			log.Println("Error reading from websocket. err:", err)
			break
		}
		if msgType != websocket.MessageText {
			log.Println("Unexpected message type")
			continue
		}
		msg, err := parseMessageRequest(payload)
		if err != nil {
			log.Println(">>>> payload:", string(payload))
			log.Println("Error parsing message. err:", err)
			continue
		}

		switch msgPayload := msg.(type) {
		case *CreateOfferRequest:
			sdpObj, err := base64.StdEncoding.DecodeString(msgPayload.OfferSDP)
			if err != nil {
				continue
			}

			marketplace.AddOffer(msgPayload.OfferID, sdpObj)
			marketplace.Display(pageLogger)
		case *AnswerOfferRequest:
			log.Println("Answer received")
			sdpObj, err := base64.StdEncoding.DecodeString(msgPayload.AnswerSDP)
			if err != nil {
				continue
			}

			peer.ReceiveOffer(sdpObj)
		default:
			log.Println("Unknown message type")
		}
	}
}

func (ws *WSClient) Writer(msg []byte) {
	err := ws.conn.Write(context.Background(), websocket.MessageText, msg)
	if err != nil {
		log.Println("Error writing to websocket. err:", err)
	}
}

func (ws *WSClient) Close() {
	ws.conn.Close(websocket.StatusNormalClosure, "")
}
