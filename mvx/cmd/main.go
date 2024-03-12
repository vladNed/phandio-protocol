package main

import (
	"log"
	// "github.com/mvx/config"
	"github.com/mvx/config"
	"github.com/mvx/contract"
)

func main() {
	log.Println("Starting MVX CLI")
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	log.Println("Configuration loaded")

	a := contract.DeploySwapContractRequest{
		TimeoutDurationOne: 1000,
		TimeoutDurationTwo: 2000,
		ClaimCommitment: "hashOne",
		RefundCommitment: "hashTwo",
		ClaimerAddress: "erd1krzz9hzh30npwuyy7vc3e6p3z7mp5rv20hnkq4xtxxnhrhtzhwsqz76cxg",
	}
	txHash, err := contract.DeploySwapContract(cfg, &a, "1000000000000000000")
	if err != nil {
		log.Fatalln("Error deploying swap contract: ", err)
		panic(err)
	}
	log.Println("Swap contract deployed. Transaction hash: ", txHash)
}
