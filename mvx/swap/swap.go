package swap

import (
	"context"
	"time"

	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/mvx/config"
)

func vmQueryRequest(callerAddress string, funcName string, config *config.Config) []byte {
	args := blockchain.ArgsProxy{
		ProxyURL:            config.ProxyURL,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}

	proxy, err := blockchain.NewProxy(args)
	if err != nil {
		panic(err)
	}

	vmRequest := &data.VmValueRequest{
		Address:    config.SwapContractAddress,
		FuncName:   funcName,
		CallerAddr: callerAddress,
		CallValue:  "",
		Args:       nil,
	}

	response, err := proxy.ExecuteVMQuery(context.Background(), vmRequest)
	if err != nil {
		panic(err)
	}

	return response.Data.ReturnData[0]
}

func GetSwapState(callerAddress string, config *config.Config) string {
	swapState := string(vmQueryRequest(callerAddress, "getState", config))
	if swapState == "" {
		return "0"
	}

	return swapState
}


func GetSecretCommitment(callerAddress string, config *config.Config) string {
	secretCommitment := string(vmQueryRequest(callerAddress, "getSecretCommitment", config))
	return secretCommitment
}