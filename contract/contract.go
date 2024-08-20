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

type ContractService struct {
	client   *ethclient.Client
	contract *bind.BoundContract
	auth     *bind.TransactOpts
}

func NewContractService() *ContractService {
	ethNodeURL := os.Getenv("ETH_NODE_URL")
	if ethNodeURL == "" {
		log.Fatalf("ETH_NODE_URL environment variable is not set")
	}

	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress("0x8464135c8f25da09e49bc8782676a84730c318bc")
	contractABI, err := utils.GetABI()
	if err != nil {
		log.Fatalf("Failed to get contract ABI: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	boundContract := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)

	return &ContractService{
		client:   client,
		contract: boundContract,
	}
}

func LoadPrivateKey() *ecdsa.PrivateKey {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	return privateKey
}

func (cs *ContractService) InitializeAuth() {
	privateKey := LoadPrivateKey()
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := cs.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := cs.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	chainID, err := cs.client.NetworkID(context.Background())
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

	cs.auth = auth
}

func (cs *ContractService) GetBalance(address common.Address) *big.Int {
	var result *big.Int
	err := cs.contract.Call(nil, &[]interface{}{&result}, "balanceOf", address)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}
	return result
}

func (cs *ContractService) MintToken(toAddress common.Address, amount *big.Int) {
	tx, err := cs.contract.Transact(cs.auth, "mint", toAddress, amount)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())
}
