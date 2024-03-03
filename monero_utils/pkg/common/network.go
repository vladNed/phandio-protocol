package common

// Environment represents the environment the swap will run in (ie. mainnet, stagenet, or development)
type Environment byte

const (
	// Undefined is a placeholder, do not pass it to functions
	Undefined Environment = iota
	// Mainnet is for real use with mainnet monero
	Mainnet
	// Stagenet is for testing with stagenet monero
	Stagenet
	// Development is for testing with a local monerod in regtest mode
	Development
)
