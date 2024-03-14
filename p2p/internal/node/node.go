package node

import (
	"bufio"
	"log"
	"os"

	"github.com/mvx-mnr-atomic/p2p/internal/peer"
	"github.com/pion/webrtc/v4"
)

type Node struct {}


func NewNode() (*Node, error) {
	return &Node{}, nil
}

func (n *Node) Start() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Println("NODE: New peer created")
		peer, err := peer.NewPeer()
		if err != nil {
			log.Println("NODE: Could not create peer. err:", err)
			return err
		}

		log.Println("Enter offer SDP: ")
		rawSDP, err := reader.ReadString('\n')
		if err != nil {
			log.Println("NODE: Could not read SDP. err:", err)
			peer.LocalConnection.Close()
			continue
		}

		err = peer.ReceiveOffer([]byte(rawSDP))
		if err != nil {
			log.Println("NODE: Could not receive offer. err:", err)
			peer.LocalConnection.Close()
			continue
		}

		err = peer.SendAnswer()
		if err != nil {
			log.Println("NODE: Could not send answer. err:", err)
			peer.LocalConnection.Close()
			continue
		}

		gatherComplete := webrtc.GatheringCompletePromise(peer.LocalConnection)
		<- gatherComplete
	}
}
