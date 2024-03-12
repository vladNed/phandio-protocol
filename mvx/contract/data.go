package contract

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/multiversx/mx-sdk-go/data"
	"log"
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
		"create_swap@%s@%s@%s@%s@%s",
		timeoutOneSerialized,
		timeoutTwoSerialized,
		claimCommitmentSerialized,
		refundCommitmentSerialized,
		claimerAddressSerialized,
	)

	return data, nil
}
