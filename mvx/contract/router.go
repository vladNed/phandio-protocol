package contract

import (
	"log"

	"github.com/mvx/config"
)

func DeploySwapContract(
	config *config.Config,
	deploySwapRequest *DeploySwapContractRequest,
	value string,
) (string, error) {
	txData, err := deploySwapRequest.Serialize()
	if err != nil {
		log.Fatalln("Error serializing deploy swap contract request: ", err)
		return "", err
	}

	txHash, err := ContractExecute(config, config.SwapRouterContractAddress, value, []byte(txData))
	if err != nil {
		log.Fatalln("Error executing contract: ", err)
		return "", err
	}

	return txHash, nil
}
