package main

import (
	"fmt"
	"os"

	"github.com/Bullrich/eth-uni-cli/wallet"
)

func main() {
	printBalances()
}

func printBalances() {
	apiKey := os.Getenv("INFURA_API_KEY")
	walletAddress := os.Getenv("WALLET_ADDRESS")

	client := wallet.NewUser(apiKey, walletAddress)

	fmt.Println(client.GetBalances())
}
