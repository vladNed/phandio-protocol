package hub

import (
	"encoding/json"
	"errors"
)

var (
	Offer  = "offer"
	Answer = "answer"
)

type Message interface{}
type MessageRequest struct {
	Type string `json:"type"`
}

func (mr *MessageRequest) Unmarshal(payload []byte) (Message, error) {
	switch mr.Type {
	case Offer:
		var req CreateOfferRequest
		return &req, unmarshalRequest(payload, &req)
	case Answer:
		var req AnswerOfferRequest
		return &req, unmarshalRequest(payload, &req)
	default:
		return nil, errors.New("invalid message type")
	}

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

type MessageResponse struct {
	Status  int         `json:"status"`
	Details interface{} `json:"details"`
}

func parseMessageRequest(payload []byte) (Message, error) {
	// Step 1: Parse the message request type
	var req MessageRequest
	err := unmarshalRequest(payload, &req)
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

func parseMessageResponse(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func unmarshalRequest(payload []byte, req interface{}) error {
	err := json.Unmarshal(payload, req)
	if err != nil {
		return err
	}
	return nil
}
