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

func GetWalletNonce(config *config.Config) (uint64, error) {
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
		return 0, err
	}

	wallet := interactors.NewWallet()
	privateKey, err := wallet.LoadPrivateKeyFromPemData([]byte(config.WalletPemData))
	if err != nil {
		log.Fatal("Error loading private key: ", err)
		return 0, err
	}

	address, err := wallet.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal("Error getting address from private key: ", err)
		return 0, err
	}

	accountInfo, err := proxy.GetAccount(context.Background(), address)
	if err != nil {
		log.Fatal("Error getting account info: ", err)
		return 0, err
	}

	return accountInfo.Nonce, nil
}

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
	txData []byte,
) (*DeploySwapContractResponse, error) {
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
		return nil, err
	}

	wallet := interactors.NewWallet()
	privateKey, err := wallet.LoadPrivateKeyFromPemData([]byte(config.WalletPemData))
	if err != nil {
		log.Fatal("Error loading private key: ", err)
		return nil, err
	}

	address, err := wallet.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Fatal("Error getting address from private key: ", err)
		return nil, err
	}

	netConfigs, err := proxy.GetNetworkConfig(context.Background())
	if err != nil {
		log.Fatalln("Error getting network config: ", err)
		return nil, err
	}

	tx, _, err := proxy.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		log.Fatal("Error getting default transaction arguments: ", err)
		return nil, err
	}
	tx.Data = txData
	tx.Receiver = contractAddress
	tx.Value = value
	tx.GasLimit = 9_000_000

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		// TODO: Log error
		return nil, err
	}

	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
	ti, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		// TODO: Log error
		return nil, err
	}
	err = ti.ApplyUserSignature(holder, &tx)
	if err != nil {
		// TODO: Log error
		return nil, err
	}
	tx.Version = 2
	tx.Options = 0
	err = ti.ApplyUserSignature(holder, &tx)
	if err != nil {
		log.Fatalln("Error applying user signature: ", err)
		return nil, nil
	}
	ti.AddTransaction(&tx)
	log.Println("Sending transaction...")
	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 10)
	if err != nil {
		log.Fatalln("Error sending transaction: ", err)
		return nil, err
	}
	log.Println("Transaction sent. Hash: ", hashes[0])

	var response *data.TransactionInfo
	poolingExitStatus := 0
	for i := 0; i < 10; i++ {
		var respErr error
		response, respErr = proxy.GetTransactionInfoWithResults(context.Background(), hashes[0])
		if respErr != nil {
			log.Fatal("Error getting transaction info with results: ", respErr)
			return nil, respErr
		}

		if response.Data.Transaction.Status == "success" {
			poolingExitStatus = 1
			break
		}

		time.Sleep(2 * time.Second)
		log.Println("Pooling tx result ....")
	}

	if poolingExitStatus == 0 {
		log.Fatal("Error: transaction could not be processed")
		return nil, err
	}

	txResponse, err := LoadFromTxOnNetwork(response.Data.Transaction)
	if err != nil {
		log.Fatal("Error loading from transaction on network: ", err)
		return nil, err
	}

	return txResponse, nil
}
