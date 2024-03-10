package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ProxyURL            string
	SwapContractAddress string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ProxyURL:            os.Getenv("PROXY_URL"),
		SwapContractAddress: os.Getenv("SWAP_CONTRACT_ADDRESS"),
	}

	return cfg, nil
}
