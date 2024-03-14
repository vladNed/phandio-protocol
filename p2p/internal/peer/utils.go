package peer

import (
	"encoding/json"
	"log"

	"github.com/pion/webrtc/v4"
)



func decodeSDP(raw_sdp []byte) (webrtc.SessionDescription, error) {
	sdp := webrtc.SessionDescription{}
	err := json.Unmarshal(raw_sdp, &sdp)
	if err != nil {
		log.Fatalln("Could not decode SDP. err:", err)
		return sdp, err
	}

	return sdp, nil
}


func encodeSDP(sdp webrtc.SessionDescription) ([]byte, error) {
	raw_sdp, err := json.Marshal(sdp)
	if err != nil {
		log.Fatalln("Could not encode SDP. err:", err)
		return nil, err
	}

	return raw_sdp, nil
}