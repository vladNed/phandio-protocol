//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"

	"encoding/base64"

	"github.com/mvx-mnr-atomic/p2p/internal/peer"
	"github.com/pion/webrtc/v4"
)

func main() {

	// Create a new peer
	peer, err := peer.NewPeer()
	if err != nil {
		fmt.Println("Could not create peer. err:", err)
		return
	}

	gatherCompletes := webrtc.GatheringCompletePromise(peer.LocalConnection)
	js.Global().Set("goCreateSDPOffer", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			offer, err := peer.CreateOffer()
			if err != nil {
				fmt.Println("Could not create offer. err:", err)
				return
			}

			offerEncoded := base64.StdEncoding.EncodeToString(offer)
			offerSDPDataElement := getElementByID("offerSDPData")
			offerSDPDataElement.Set("value", offerEncoded)
			fmt.Println("goCreateSDPOffer... DONE")
		}()
		return js.Undefined()
	}))

	js.Global().Set("goCreateSDPAnswer", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			// TODO: Add code here
		}()

		return js.Undefined()
	}))
	<-gatherCompletes

	// NOTE: This is a blocking call
	// This is a temporary solution to keep the program running.
	// Should be done more elegantly like having an event loop
	// and handling events in a non-blocking manner.
	select {}
}

func getElementByID(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}
