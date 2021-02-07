package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Bullrich/eth-uni-cli/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	startWebServer()
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

	app.Get("/api/balance/:address", func(c *fiber.Ctx) error {
		user := wallet.NewUser(apiKey, c.Params("address"))
		balance := user.GetBalances()
		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(string(balanceJSON))
	})

	app.Get("/api/tx/", func(c *fiber.Ctx) error {
		address := c.Query("address")
		coin := c.Query("coin")
		input := c.Query("input-amount")
		output := c.Query("output-amount")
		privateKey := c.Query("private-key")

		inputValue, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return c.SendString("Invalid input. Must be a number")
		}

		outputValue, err := strconv.ParseInt(output, 10, 64)
		if err != nil {
			return c.SendString("Invalid output. Must be a number")
		}

		client := wallet.NewUser(apiKey, address)

		transaction, err := wallet.BuyCoin(client, privateKey, coin, inputValue, outputValue)
		if err != nil {
			return c.SendString(fmt.Sprintf("%+v", err))
		}

		response := &fiber.Map{
			"transactionHash": transaction,
			"message":         fmt.Sprintf("tx sent. See the transaction in: https://etherscan.io/tx/%s\n", transaction),
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(string(responseJSON))
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
