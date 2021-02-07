# Ethereum Uniswap CLI

## How to

You need to get an api key from [infura.io](https://infura.io/)

### With Go installed in your system

Clone the repo and execute the following:
```shell
go get .
INFURA_API_KEY=<infura-api-key> go run main.go
```

### With docker

```shell
docker build -t eth-uni-cli .
docker run -dp 3000:3000 -e INFURA_API_KEY='<infura-api-key>' eth-uni-cli
```

Then go to http://127.0.0.1:3000/balance to see the balance of a wallets in a simple UI

## API

- `/api/balance/<address>`: Returns the balance in `dai`, `sai`, `wei` and `mkr` of the given address
- `/api/tx`: Exchanges a ether for a given token. See the parameters to know how to correctly use it.
  - To use this endpoint, you need to pass the parameters in a get call: `/api/tx?address=<address>&coin=<dai|sai|mkr>&input-amount=<int>&output-amount=<int>&private-key=<key>`

## Scope

Build an API to interact with Ethereum and Uniswap.
One endpoint should get the balances of `DAI`, `MKR` and one `Token of your choice` for a specific ethereum address.
Another endpoint should make a Uniswap swap with the pairs `ETH/DAI`, `ETH/MKR` and `ETH/{Token of your choice}`.

Please work on this project for a fixed amount of time (recommended 4h-6h). The aim is not to complete this task in the fixed amout of time, but rather show the progress you made until then.

Recommended is to use the [go-ethereum](https://github.com/ethereum/go-ethereum) library and a webserver framework like [echo](https://github.com/labstack/echo) or [fiber](https://github.com/gofiber/fiber).

Also it's recommended to connect to a ethereum testnet node to save on transaction fees while developing.

If you feel like not using any of the above, feel free to choose the technology you are feeling most comfortable with. This is just for reference.
Also, if you choose to use another programming language, its free of choice. (Just dont use javascript, typescript or nodejs)

## Resource

- Rinkeby network address: 0x2dDF0C1A65ABddeF67796b5fd396F1bc9D57C6EC
- Private Key for this account with Rinkeby ETH: via PM
- Example Rinkeby Pool Contract [0x7a250d5630b4cf539739df2c5dacb4c659f2488d](https://rinkeby.etherscan.io/address/0x7a250d5630b4cf539739df2c5dacb4c659f2488d#writeContract)

## Setup

- Fork this repository
- Develop on a feature branch and push your changes frequently
- Finally create a pull request to this repository
