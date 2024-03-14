package peer

import (
	"log"
	"os"

	"github.com/pion/webrtc/v4"
)

type Peer struct {
	LocalConnection *webrtc.PeerConnection
	DataChannel     *webrtc.DataChannel
	state           PeerState
}

// / NewPeer creates a new webRTC peer connection with a dedicated
// / data channel.
func NewPeer() (*Peer, error) {
	rtcCfg := GetICEConfiguration()
	peerConnection, err := webrtc.NewPeerConnection(rtcCfg)
	if err != nil {
		log.Fatalln("Could not create peer. err:", err)
		return nil, err
	}

	peerDataChannel, err := peerConnection.CreateDataChannel(DATA_CHANNEL_LABEL, nil)
	if err != nil {
		log.Fatalln("Could not create data channel. err:", err)
		return nil, err
	}

	peer := &Peer{
		LocalConnection: peerConnection,
		DataChannel:     peerDataChannel,
		state:           PeerStateDefault,
	}

	peer.setupConnectionCallbacks()
	peer.setupDataChannelProtocol()

	return peer, nil
}


func (p *Peer) setupConnectionCallbacks() {
	p.LocalConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		log.Println("Peer connection state has changed: ", connectionState)

		switch connectionState {
		case webrtc.PeerConnectionStateConnecting:
			log.Println("Connection state. Initial Peer state 0")
			p.state = PeerState0
		case webrtc.PeerConnectionStateConnected:
			if p.state != PeerState0 {
				log.Fatalln("Cannot connect to peer. Invalid state")
				os.Exit(1)
			}
			log.Println("Connection established. Peer state is now 1.")
			p.state = PeerState1
		case webrtc.PeerConnectionStateDisconnected:
			log.Println("Peer connection state is disconnected")
			p.state = PeerStateDefault
		case webrtc.PeerConnectionStateFailed:
			log.Println("Peer connection state has failed")
		case webrtc.PeerConnectionStateClosed:
			log.Println("Peer connection state is closed")
		}
	})
	p.LocalConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		log.Println("New ICE candidate found")
	})
}

func(p *Peer) onInvalidDataChannel(d *webrtc.DataChannel) {
	log.Println("Ignoring data channel. Invalid label")
	err := d.Close()
	if err != nil {
		log.Fatalln("Could not close data channel. err:", err)
		os.Exit(1)
	}

	err = p.LocalConnection.Close()
	if err != nil {
		log.Fatalln("Could not close peer connection. err:", err)
		os.Exit(1)
	}
}

func (p *Peer) setupDataChannelProtocol() {
	// Set data channel protocol
	p.DataChannel.OnOpen(func() {
		log.Println("Data channel is open")

		candidatePair, err := p.LocalConnection.SCTP().Transport().ICETransport().GetSelectedCandidatePair()
		if err != nil {
			log.Fatalln("Could not get selected candidate pair. err:", err)
			return
		}

		log.Println("Selected candidate pair: ", candidatePair)
		// TODO: Add on open channel protocol logic
	})
	p.DataChannel.OnClose(func() {
		log.Println("Data channel is closed")

		// TODO: Add on close channel protocol logic
	})
	p.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Printf("Message received: %s\n", string(msg.Data))

		// TODO: Add message handling logic
	})

	// Set receiving data channel
	p.LocalConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		log.Println("Data-Channel established")

		if d.Label() != DATA_CHANNEL_LABEL {
			p.onInvalidDataChannel(d)
			return
		}

		d.OnOpen(func () {
			log.Printf("Accepted data channel. %s - %d\n", d.Label(), d.ID())
			// TODO: Add authentication logic
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Message received: %s\n", string(msg.Data))
			// TODO: Add message handling logic
		})
	})
}


func (p *Peer) ReceiveOffer(offer []byte) error {
	offer_sdp, err := decodeSDP(offer)
	if err != nil {
		log.Fatalln("Could not decode offer. err:", err)
		return err
	}

	err = p.LocalConnection.SetRemoteDescription(offer_sdp)
	if err != nil {
		log.Fatalln("Could not set remote description. err:", err)
		return err
	}

	return nil
}


func (p *Peer) SendAnswer() error {
	answerSDP, err := p.LocalConnection.CreateAnswer(nil)
	if err != nil {
		log.Fatalln("Could not create answer. err:", err)
		return err
	}

	answer, err := encodeSDP(answerSDP)
	if err != nil {
		log.Fatalln("Could not encode answer. err:", err)
		return err
	}

	log.Println("Sending answer")
	err = p.LocalConnection.SetLocalDescription(answerSDP)
	if err != nil {
		log.Fatalln("Could not set local description. err:", err)
		return err
	}

	// TODO: Implement signalling logic
	log.Println("Answer ->>> ", string(answer))

	return nil
}

func (p *Peer) CreateOffer() ([]byte, error) {
	offer, err := p.LocalConnection.CreateOffer(nil)
	if err != nil {
		log.Fatalln("Could not create offer. err:", err)
		return nil, err
	}

	err = p.LocalConnection.SetLocalDescription(offer)
	if err != nil {
		log.Fatalln("Could not set local description. err:", err)
		return nil, err
	}

	offerSDP, err := encodeSDP(offer)
	if err != nil {
		log.Fatalln("Could not encode offer. err:", err)
		return nil, err
	}

	return offerSDP, nil

}
