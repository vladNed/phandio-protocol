package signalling

import (
	"encoding/json"
	"errors"
)

var (
	Offer  = "offer"
	Answer = "answer"
)

// TODO: All schemas json messages should be taken from signalling package
type Message interface{}

type MessageRequest struct {
	Type string `json:"type"`
}

func (mr *MessageRequest) Unmarshal(payload []byte) (Message, error) {
	switch mr.Type {
	case Offer:
		var req CreateOfferRequest
		return &req, json.Unmarshal(payload, &req)
	case Answer:
		var req AnswerOfferRequest
		return &req, json.Unmarshal(payload, &req)
	default:
		return nil, errors.New("invalid message type")
	}

}

type RegisterChannelRequest struct {
	Channels []string `json:"channel"`
}

type Response struct {
	Status  int         `json:"status"`
	Details interface{} `json:"details"`
}


func parseMessageRequest(payload []byte) (Message, error) {
	// Step 1: Parse the message request type
	var req MessageRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	// Step 2: Parse the message fully
	requestMessage, err := req.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	return requestMessage, nil
}


func parseResponse(payload []byte) (*Response, error) {
	var res Response
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type CreateOfferRequest struct {
	Type     string `json:"type"`
	OfferID  string `json:"offerId"`
	OfferSDP string `json:"sdp"`
}

type AnswerOfferRequest struct {
	Type      string `json:"type"`
	OfferID   string `json:"offerId"`
	AnswerSDP string `json:"sdp"`
}
