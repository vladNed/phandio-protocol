package main

import (
	"log"
	// "github.com/mvx/config"
	"github.com/mvx/contract"
)

func main() {
	log.Println("Starting MVX CLI")
	// cfg, err := config.LoadConfig()
	// if err != nil {
	// 	panic(err)
	// }
	log.Println("Configuration loaded")

	a := contract.DeploySwapContractRequest{
		TimeoutDurationOne: 100,
		TimeoutDurationTwo: 200,
		ClaimCommitment: "hashOne",
		RefundCommitment: "hashTwo",
		ClaimerAddress: "erd1qqqqqqqqqqqqqpgqftd63554n7e89t6ek7k8n5aweydhcsh0hwsq24k362",
	}
	serializedData, err := a.Serialize()
	if err != nil {
		log.Fatalln("Error serializing data: ", err)
	}
	log.Println(serializedData)

	// TODO: Add your code here for CLI
}
