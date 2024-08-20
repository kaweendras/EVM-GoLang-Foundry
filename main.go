package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kaweendras/EVM-GoLang-Foundry/utils"
)

func main() {
	ethNodeURL := os.Getenv("ETH_NODE_URL")
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	// Connect to an Ethereum node
	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	contractAddress := common.HexToAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3")
	contractABI, err := utils.GetABI()
	if err != nil {
		log.Fatalf("Failed to get contract ABI: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contract := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)
	toAddress := common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720")

	//print the token amount of an address (READ)
	var result *big.Int
	err = contract.Call(nil, &[]interface{}{&result}, "balanceOf", toAddress)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}
	fmt.Printf("Token amount: %s\n", result.String())

	//mint token to an address (WRITE)
	// amount := big.NewInt(1000000000000000000)
	// tx, err := contract.Transact(auth, "mint", toAddress, amount)
	// if err != nil {
	// 	log.Fatalf("Failed to send transaction: %v", err)
	// }
	// fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	//burn token from an address (WRITE)
	// amount := big.NewInt(1500000000000000000)
	// tx, err := contract.Transact(auth, "burnFrom", toAddress, amount)
	// if err != nil {
	// 	log.Fatalf("Failed to send transaction: %v", err)
	// }
	// fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

}
