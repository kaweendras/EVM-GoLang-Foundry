package client

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Connect() *ethclient.Client {
	ethNodeURL := os.Getenv("ETH_NODE_URL")
	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	return client
}
