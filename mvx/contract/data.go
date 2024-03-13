package contract

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/data"
)

type DeploySwapContractRequest struct {
	TimeoutDurationOne uint64
	TimeoutDurationTwo uint64
	ClaimCommitment    string
	RefundCommitment   string
	ClaimerAddress     string
}

// / Serialize the deploy swap contract request to the format expected to be sent in the
// / transaction data field.
func (r *DeploySwapContractRequest) Serialize() (string, error) {

	// Serialize the timeout duration one
	timeoutOneBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(timeoutOneBuf, r.TimeoutDurationOne)
	timeoutOneSerialized := fmt.Sprintf("%X", timeoutOneBuf)

	// Serialize the timeout duration two
	timeoutTwoBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(timeoutTwoBuf, r.TimeoutDurationTwo)
	timeoutTwoSerialized := fmt.Sprintf("%X", timeoutTwoBuf)

	// Serialize the claim commitment
	claimCommitmentBytes := []byte(r.ClaimCommitment)
	claimCommitmentSerialized := hex.EncodeToString(claimCommitmentBytes)

	// Serialize the refund commitment
	refundCommitmentBytes := []byte(r.RefundCommitment)
	refundCommitmentSerialized := hex.EncodeToString(refundCommitmentBytes)

	// Serialize the claimer address
	address, err := data.NewAddressFromBech32String(r.ClaimerAddress)
	if err != nil {
		log.Fatalln("Invalid claimer address provided: ", err)
		return "", err
	}
	address.AddressBytes()
	claimerAddressBytes := address.AddressBytes()
	claimerAddressSerialized := hex.EncodeToString(claimerAddressBytes)

	data := fmt.Sprintf(
		"createSwap@%s@%s@%s@%s@%s",
		timeoutOneSerialized,
		timeoutTwoSerialized,
		claimCommitmentSerialized,
		refundCommitmentSerialized,
		claimerAddressSerialized,
	)

	return data, nil
}

type DeploySwapContractResponse struct {
	respAddress string
	status      string
}

func LoadFromTxOnNetwork(txData data.TransactionOnNetwork) (*DeploySwapContractResponse, error) {
	status := getDeploySwapTxStatus(txData.Logs.Events)
	if status != "successful" {
		return nil, errors.New("transaction failed")
	}

	scResult, err := deserializeDeploySwapScResults(string(txData.ScResults[0].Data))
	if err != nil {
		log.Fatalln("Error deserializing deploy swap contract results: ", err)
		return nil, err
	}

	return &DeploySwapContractResponse{
		respAddress: scResult[1],
		status:      status,
	}, nil
}

func deserializeDeploySwapScResults(txData string) ([]string, error) {
	parts := strings.Split(txData, "@")
	results := []string{}

	// Deserialize the status
	status, err := hex.DecodeString(parts[1])
	if err != nil {
		log.Fatalln("Error deserializing status: ", err)
		return nil, err
	}
	results = append(results, string(status))

	// Deserialize the swap contract address
	swapContractAddress, err := hex.DecodeString(parts[2])
	if err != nil {
		log.Fatalln("Error deserializing swap contract address: ", err)
		return nil, err
	}
	address := data.NewAddressFromBytes(swapContractAddress)
	addressString, err := address.AddressAsBech32String()
	if err != nil {
		log.Fatalln("Error converting swap contract address to bech32: ", err)
		return nil, err
	}
	results = append(results, addressString)


	return results, nil
}

func getDeploySwapTxStatus(txLogEvents []*transaction.Events) string {
	for _, event := range txLogEvents {
		if string(event.Identifier) == "completedTxEvent" {
			return "successful"
		}
	}

	return "failed"
}
