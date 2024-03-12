package contract

import (
	"context"
	"log"
	"time"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/mvx/config"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

func ContractQuery(
	config *config.Config,
	contractAddress string,
	callerAddress string,
	funcName string,
) ([]byte, error) {
	log.Println("Querying contract. Function name: ", funcName)
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
		return nil, err
	}
	vmRequest := &data.VmValueRequest{
		Address:    contractAddress,
		FuncName:   funcName,
		CallerAddr: callerAddress,
		CallValue:  "",
		Args:       nil,
	}
	response, err := proxy.ExecuteVMQuery(context.Background(), vmRequest)
	if err != nil {
		log.Fatalln("Error executing VM query: ", err)
		return nil, err
	}

	return response.Data.ReturnData[0], nil
}

func ContractExecute(
	config *config.Config,
	contractAddress string,
	value string,
	data []byte,
) {
	log.Println("Calling contract. Contract address: ", contractAddress)
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
		log.Fatal("Error creating proxy: ", err)
		return
	}

	wallet := interactors.NewWallet()
	privateKey, err := wallet.LoadPrivateKeyFromPemData([]byte(config.WalletPemData))
	if err != nil {
		log.Fatal("Error loading private key: ", err)
		return
	}

	address, err := wallet.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal("Error getting address from private key: ", err)
		return
	}

	netConfigs, err := proxy.GetNetworkConfig(context.Background())
	if err != nil {
		log.Fatalln("Error getting network config: ", err)
		return
	}

	tx, _, err := proxy.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		log.Fatal("Error getting default transaction arguments: ", err)
		return
	}
	tx.Data = data
	tx.Receiver = contractAddress
	tx.Value = value
	tx.GasLimit = 6000000

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		// TODO: Log error
		return
	}

	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
	ti, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		// TODO: Log error
		return
	}
	err = ti.ApplyUserSignature(holder, &tx)
	if err != nil {
		// TODO: Log error
		return
	}
	ti.AddTransaction(&tx)

	// a new transaction with the signature done on the hash of the transaction
	// it's ok to reuse the arguments here, they will be copied, anyway
	tx.Version = 2
	tx.Options = 0
	tx.Nonce++ // do not forget to increment the nonce, otherwise you will get 2 transactions
	// with the same nonce (only one of them will get executed)
	err = ti.ApplyUserSignature(holder, &tx)
	if err != nil {
		log.Fatalln("Error applying user signature: ", err)
		return
	}
	ti.AddTransaction(&tx)
	log.Println("Sending transaction...")
	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 10)
	if err != nil {
		log.Fatalln("Error sending transaction: ", err)
		return
	}

	log.Println("Transaction hashes: ", hashes)
}
