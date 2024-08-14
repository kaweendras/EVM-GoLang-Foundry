package contract

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

func LoadPrivateKey() *ecdsa.PrivateKey {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	return privateKey
}

func GetAuth(client *ethclient.Client, privateKey *ecdsa.PrivateKey) *bind.TransactOpts {
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

	return auth
}

func GetContract(client *ethclient.Client) *bind.BoundContract {
	contractAddress := common.HexToAddress("0x8464135c8f25da09e49bc8782676a84730c318bc")
	contractABI, err := utils.GetABI()
	if err != nil {
		log.Fatalf("Failed to get contract ABI: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	return bind.NewBoundContract(contractAddress, parsedABI, client, client, client)
}

func GetBalance(contract *bind.BoundContract, address common.Address) *big.Int {
	var result *big.Int
	err := contract.Call(nil, &[]interface{}{&result}, "balanceOf", address)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}
	return result
}

func MintToken(auth *bind.TransactOpts, contract *bind.BoundContract, toAddress common.Address, amount *big.Int) {
	tx, err := contract.Transact(auth, "mint", toAddress, amount)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())
}
