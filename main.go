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
)

func main() {
	// Connect to an Ethereum node
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
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

	contractAddress := common.HexToAddress("0x8464135c8f25da09e49bc8782676a84730c318bc")
	contractABI, err := getABI()
	if err != nil {
		log.Fatalf("Failed to get contract ABI: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contract := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)

	// Interact with the contract (Write)
	// tx, err := contract.Transact(auth, "setNumber", big.NewInt(1000))
	// if err != nil {
	// 	log.Fatalf("Failed to send transaction: %v", err)
	// }
	// fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	//print the value of the number(Read)
	var number *big.Int
	var numberSlice []interface{}
	err = contract.Call(&bind.CallOpts{Context: context.Background()}, &numberSlice, "number")
	number = numberSlice[0].(*big.Int)
	if err != nil {
		log.Fatalf("Failed to call contract function: %v", err)
	}

	fmt.Printf("Variable Value: %s\n", number.String())

}

func getABI() ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %v", err)
	}
	fmt.Println("Current working directory:", cwd)

	abiFile, err := os.ReadFile("ABI/bamla.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %v", err)
	}

	return abiFile, nil
}
