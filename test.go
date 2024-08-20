package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kaweendras/EVM-GoLang-Foundry/contract"
)

func test() {
	privateKey := contract.LoadPrivateKey()
	contractService := contract.NewContractService()
	contractService.InitializeAuth(privateKey)

	toAddress := common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720")

	// Print the token amount of an address (READ)
	balance := contractService.GetBalance(toAddress)
	fmt.Printf("Token amount: %s\n", balance.String())

	// Mint token to an address (WRITE)
	amount := big.NewInt(1000000000000000000)
	contractService.MintToken(toAddress, amount)
}
