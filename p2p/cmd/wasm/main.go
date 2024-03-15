//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"log"
	"syscall/js"

	"encoding/base64"
	"encoding/hex"

	"github.com/mvx-mnr-atomic/p2p/internal/peer"
	"github.com/mvx-mnr-atomic/p2p/internal/monero/crypto"
	"github.com/mvx-mnr-atomic/p2p/internal/monero/common"
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
