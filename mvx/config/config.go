package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ProxyURL                  string
	SwapContractAddress       string
	SwapRouterContractAddress string
	WalletPemData             string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	pemFile, err := os.ReadFile(os.Getenv("WALLET_PEM_FILE_PATH"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ProxyURL:                  os.Getenv("PROXY_URL"),
		SwapContractAddress:       os.Getenv("SWAP_CONTRACT_ADDRESS"),
		SwapRouterContractAddress: os.Getenv("SWAP_ROUTER_CONTRACT_ADDRESS"),
		WalletPemData:             string(pemFile),
	}

	return cfg, nil
}
