package main

import (
	"github.com/mvx/swap"
	"github.com/mvx/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	// TODO: Add your code here for CLI
}
