//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"log"
	"syscall/js"

	"encoding/base64"

	"github.com/mvx-mnr-atomic/p2p/internal/peer"
)

func main() {

	localPeer, err := peer.NewPeer(logger)
	if err != nil {
		log.Println("Could not create peer. err:", err)
		return
	}
	logger("Peer initialized.")
	js.Global().Set("goCreateSDPOffer", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			newSDPOffer, err := localPeer.CreateOffer()
			if err != nil {
				logger("Could not create offer. err:" + err.Error())
				return
			}
			offerEncoded := base64.StdEncoding.EncodeToString(newSDPOffer)
			getElementByID("offerSDPData").Set("value", offerEncoded)
		}()
		return js.Undefined()
	}))
	js.Global().Set("goCreateSDPAnswer", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			offerSDPRaw := getElementByID("offerSDP").Get("value").String()
			offerSDPJson, err := base64.StdEncoding.DecodeString(offerSDPRaw)
			if err != nil {
				logger("Could not decode offer. err:" + err.Error())
				return
			}
			err = localPeer.ReceiveOffer(offerSDPJson)
			if err != nil {
				logger("Could not parse offer. err:" + err.Error())
				return
			}
			answerStr, err := localPeer.SendAnswer()
			if err != nil {
				logger("Could not create answer. err:" + err.Error())
				return
			}
			answerEncoded := base64.StdEncoding.EncodeToString(answerStr)
			getElementByID("answerSDPData").Set("value", answerEncoded)
		}()
		return js.Undefined()
	}))
	js.Global().Set("goSetSDPAnswer", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			if localPeer.LocalConnection.RemoteDescription() != nil {
				logger("P2P Connection already established.")
				return
			}
			answerSDPRaw := getElementByID("peerID").Get("value").String()
			answerSDP, err := base64.StdEncoding.DecodeString(answerSDPRaw)
			if err != nil {
				fmt.Println("Could not decode answer. err:", err)
				return
			}
			err = localPeer.ReceiveOffer(answerSDP)
			if err != nil {
				fmt.Println("Could not parse answer. err:", err)
				return
			}
			logger("P2P Connection established.")
		}()
		return js.Undefined()
	}))
	js.Global().Set("goSendData", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			msg := getElementByID("sendData").Get("value").String()
			localPeer.DataChannel.Send([]byte(msg))
		}()
		return js.Undefined()
	}))

	// NOTE: This is a blocking call
	// This is a temporary solution to keep the program running.
	// Should be done more elegantly like having an event loop
	// and handling events in a non-blocking manner.
	select {}
}

func getElementByID(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func logger(msg string) {
	el := getElementByID("logs")
	el.Set("innerHTML", el.Get("innerHTML").String()+"> "+msg+"<br>")
}
