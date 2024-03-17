package hub

import (
	"encoding/json"
	"errors"
)

var (
	MarketplaceChannel = "marketplace"
	OffersChannel      = "offers"
	AllChannels        = []string{MarketplaceChannel, OffersChannel}
)

type ChannelRequest struct {
	Channels []string `json:"channel"`
}

func parseChannelRequest(payload []byte) (*ChannelRequest, error) {
	var req ChannelRequest
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}
	if len(req.Channels) == 0 {
		return nil, errors.New("no channels provided")
	}
	if !verifyChannel(&req) {
		return nil, errors.New("invalid channel request")
	}
	return &req, nil
}

func containsChannel(channels []string, channel string) bool {
	for _, c := range channels {
		if c == channel {
			return true
		}
	}
	return false
}

func verifyChannel(req *ChannelRequest) bool {
	for _, channel := range req.Channels {
		if containsChannel(AllChannels, channel) {
			continue
		}
		return false
	}
	return true
}
