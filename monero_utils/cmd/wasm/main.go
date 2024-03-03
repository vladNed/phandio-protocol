//go:build js && wasm
// +build js,wasm

package main

import (
	"bbogdan95/moneroutils/pkg/common"
	mcrypto "bbogdan95/moneroutils/pkg/crypto/monero"
	"encoding/hex"
	"syscall/js"
)

func generateWalletCmd(this js.Value, p []js.Value) interface{} {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			privateKeys, err := mcrypto.GenerateKeys()
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

			pvk, err := mcrypto.NewPrivateSpendKey(pkBytes)
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

			pvk, err := mcrypto.NewPrivateSpendKey(pkBytes)
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
			pvkA, err := mcrypto.NewPrivateSpendKey(pkABytes)
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
			pvkB, err := mcrypto.NewPrivateSpendKey(pkBBytes)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			sum := mcrypto.SumPrivateSpendKeys(pvkA, pvkB)
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

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("add", js.FuncOf(concatFunction))

	js.Global().Set("generateWallet", js.FuncOf(generateWalletCmd))
	js.Global().Set("generatePubFromPriv", js.FuncOf(generatePubFromPriv))
	js.Global().Set("generateViewFromSpend", js.FuncOf(generateViewFromSpend))
	js.Global().Set("sumPrivateSpendKeys", js.FuncOf(sumPrivateSpendKeys))

	<-c
}
