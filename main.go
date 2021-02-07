package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Bullrich/eth-uni-cli/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

const uniswapRouterAddress = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"

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
	engine := html.New("./views", ".html")
	engine.Delims("{{", "}}")

	apiKey := os.Getenv("INFURA_API_KEY")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/balance/", func(c *fiber.Ctx) error {
		return c.Render("balances", nil)
	})

	app.Get("/balance/:address", func(c *fiber.Ctx) error {
		address := c.Params("address")
		balanceMap := obtainFormattedBalance(apiKey, address)

		return c.Render("walletBalance", balanceMap)
	})

	app.Post("/balance/:address", func(c *fiber.Ctx) error {
		user := wallet.NewUser(apiKey, c.Params("address"))
		balance := user.GetBalances()
		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(string(balanceJSON))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func obtainFormattedBalance(apiKey string, address string) *fiber.Map {
	user := wallet.NewUser(apiKey, address)
	if user == nil {
		return &fiber.Map{
			"validAddress": false,
		}
	}

	balance := user.GetBalances()

	return &fiber.Map{
		"address":      address,
		"wei":          balance["wei"],
		"sai":          balance["sai"],
		"mkr":          balance["mkr"],
		"dai":          balance["dai"],
		"validAddress": true,
	}
}
