package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bullrich/eth-uni-cli/wallet"
	"github.com/gofiber/fiber/v2"
)

func main() {
	startWebServer()
}

func printBalances() {
	apiKey := os.Getenv("INFURA_API_KEY")
	walletAddress := os.Getenv("WALLET_ADDRESS")

	client := wallet.NewUser(apiKey, walletAddress)

	fmt.Println(client.GetBalances())
}

func startWebServer() {
	app := fiber.New()

	apiKey := os.Getenv("INFURA_API_KEY")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/balance/:address", func(c *fiber.Ctx) error {
		user := wallet.NewUser(apiKey, c.Params("address"))
		balance := user.GetBalances()
		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(string(balanceJSON))
	})

	app.Listen(":3000")
}
