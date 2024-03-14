package main

import (
	"github.com/mvx-mnr-atomic/p2p/internal/node"
)

func main() {
	node, err := node.NewNode()
	if err != nil {
		panic(err)
	}
	node.Start()
}
