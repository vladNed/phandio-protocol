//go:build js && wasm
// +build js,wasm

package main

import (
	"log"
	"syscall/js"
	"time"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mvx-mnr-atomic/p2p/internal/monero/common"
	"github.com/mvx-mnr-atomic/p2p/internal/monero/crypto"
	"github.com/mvx-mnr-atomic/p2p/internal/peer"
	"github.com/mvx-mnr-atomic/p2p/internal/signalling"
)

func main() {
	localPeer, wsConn, marketplace := initialize()
	go wsConn.Listener(logger, localPeer, marketplace)
	js.Global().Set("goCreateOffer", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go createOffer(localPeer, wsConn)
		return js.Undefined()
	}))
	js.Global().Set("goAnswerOffer", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go answerOffer(marketplace, localPeer, wsConn)
		return js.Undefined()
	}))
	js.Global().Set("goSendData", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			msg := getElementByID("sendData").Get("value").String()
			localPeer.DataChannel.Send([]byte(msg))
		}()
		return js.Undefined()
	}))
	js.Global().Set("add", js.FuncOf(concatFunction))
	js.Global().Set("generateWallet", js.FuncOf(generateWalletCmd))
	js.Global().Set("generatePubFromPriv", js.FuncOf(generatePubFromPriv))
	js.Global().Set("generateViewFromSpend", js.FuncOf(generateViewFromSpend))
	js.Global().Set("sumPrivateSpendKeys", js.FuncOf(sumPrivateSpendKeys))

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

func generateWalletCmd(this js.Value, p []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			privateKeys, err := crypto.GenerateKeys()
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			publicKeys := privateKeys.PublicKeyPair()

			data := map[string]interface{}{
				"privsk":      privateKeys.SpendKey().Hex(),
				"pubsk":       publicKeys.SpendKey().Hex(),
				"privvk":      privateKeys.ViewKey().Hex(),
				"pubvk":       publicKeys.SpendKey().Hex(),
				"mainnet":     publicKeys.Address(common.Mainnet).String(),
				"stagenet":    publicKeys.Address(common.Stagenet).String(),
				"development": publicKeys.Address(common.Development).String(),
			}

			resolve.Invoke(js.ValueOf(data))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func generatePubFromPriv(this js.Value, p []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			pkBytes, err := hex.DecodeString(p[0].String())
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			pvk, err := crypto.NewPrivateSpendKey(pkBytes)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			data := map[string]interface{}{
				"priv": pvk.Hex(),
				"pub":  pvk.Public().Hex(),
			}

			resolve.Invoke(js.ValueOf(data))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func generateViewFromSpend(this js.Value, p []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			pkBytes, err := hex.DecodeString(p[0].String())
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			pvk, err := crypto.NewPrivateSpendKey(pkBytes)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			keyPair, err := pvk.AsPrivateKeyPair()
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			data := map[string]interface{}{
				"privsk": keyPair.SpendKey().Hex(),
				"privvk": keyPair.ViewKey().Hex(),
			}

			resolve.Invoke(js.ValueOf(data))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func sumPrivateSpendKeys(this js.Value, p []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			pkABytes, err := hex.DecodeString(p[0].String())
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}
			pvkA, err := crypto.NewPrivateSpendKey(pkABytes)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			pkBBytes, err := hex.DecodeString(p[1].String())
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}
			pvkB, err := crypto.NewPrivateSpendKey(pkBBytes)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			sum := crypto.SumPrivateSpendKeys(pvkA, pvkB)
			sumkp, err := sum.AsPrivateKeyPair()
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			data := map[string]interface{}{
				"privsk":      sumkp.SpendKey().Hex(),
				"pubsk":       sumkp.PublicKeyPair().SpendKey().Hex(),
				"privvk":      sumkp.ViewKey().Hex(),
				"pubvk":       sumkp.PublicKeyPair().SpendKey().Hex(),
				"mainnet":     sumkp.PublicKeyPair().Address(common.Mainnet).String(),
				"stagenet":    sumkp.PublicKeyPair().Address(common.Stagenet).String(),
				"development": sumkp.PublicKeyPair().Address(common.Development).String(),
			}

			resolve.Invoke(js.ValueOf(data))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func concatFunction(this js.Value, p []js.Value) interface{} {
	sum := p[0].String() + p[1].String()
	return js.ValueOf(sum)
}

// P2P Functions

// This method initialized all resources needed for the P2P node and server
// communication.
func initialize() (*peer.Peer, *signalling.WSClient, *signalling.Marketplace) {
	logger("Starting node...")
	for {
		marketplace := signalling.NewMarketplace()
		localPeer, err := peer.NewPeer(logger)
		if err != nil {
			retryDelay()
			continue
		}

		wsConn, err := signalling.NewWSClient()
		if err != nil {
			retryDelay()
			continue
		}
		res := wsConn.RegisterChannels()
		if !res {
			wsConn.Close()
			retryDelay()
			continue
		}

		logger("Node up and running.")
		return localPeer, wsConn, marketplace
	}
}

// This method is used to retry the initialization of the node in case of failure.
func retryDelay() {
	retryDelay := 5 * time.Second
	log.Println("Retrying initializing node")
	time.Sleep(retryDelay)
}

// Creates an SDP offer and sends it to the signalling server all other clients
// registered to the offers channel will receive this offer including the peer.
func createOffer(peer *peer.Peer, wsConn *signalling.WSClient) {
	newSDPOffer, err := peer.CreateOffer()
	if err != nil {
		logger("Could not create offer. err:" + err.Error())
		return
	}
	offerEncoded := base64.StdEncoding.EncodeToString(newSDPOffer)
	offerMessage := &signalling.CreateOfferRequest{
		Type: signalling.Offer,
		OfferSDP: offerEncoded,
		OfferID: uuid.NewString(),
	}
	msgPayload, err := json.Marshal(offerMessage)
	if err != nil {
		logger("Could not marshal offer. err:" + err.Error())
		return
	}

	wsConn.Writer(msgPayload)
	logger("Offer sent.")
}

// This method is used to answer an offer received from the signalling server.
// The answer is sent to the signalling server and the peer that sent the offer
// will receive the answer.
//
// All querying is done based on the offer id.
func answerOffer(marketplace *signalling.Marketplace, localPeer *peer.Peer, wsConn *signalling.WSClient) {
	offerId := getElementByID("offerId").Get("value").String()
	offerSDP, ok := marketplace.GetOffer(offerId)
	if !ok {
		logger("Offer not found.")
		return
	}
	err := localPeer.ReceiveOffer(offerSDP)
	if err != nil {
		logger("Invalid offer provided.")
		return
	}
	answerSDP, err := localPeer.SendAnswer()
	if err != nil {
		logger("Could not create answer. err:" + err.Error())
		return
	}
	answerEncoded := base64.StdEncoding.EncodeToString(answerSDP)
	answerMessage := &signalling.AnswerOfferRequest{
		Type: signalling.Answer,
		OfferID: offerId,
		AnswerSDP: answerEncoded,
	}
	msgPayload, err := json.Marshal(answerMessage)
	if err != nil {
		logger("Could not marshal answer. err:" + err.Error())
		return
	}
	wsConn.Writer(msgPayload)
	logger("Answer sent.")
}
