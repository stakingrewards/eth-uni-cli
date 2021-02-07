package wallet

import (
	"context"
	"fmt"
	"math/big"

	token "github.com/Bullrich/eth-uni-cli/token"
	"github.com/Bullrich/eth-uni-cli/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// User is a container of eht client and the address of the wallet
type User struct {
	client  *ethclient.Client
	address common.Address
}

const daiContract = "0x6B175474E89094C44Da98b954EedeAC495271d0F"
const mkrContract = "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2"
const saiContract = "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"

// NewUser constructs a User object to interact with the ethereum network
func NewUser(infuraAPIKey string, address string) *User {
	if !(utils.IsAddressValid(address)) {
		return nil
	}

	infuraAddress := fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraAPIKey)

	client, err := ethclient.Dial(infuraAddress)
	utils.CheckError(err)

	account := common.HexToAddress(address)

	return &User{client: client, address: account}
}

// GetWeiBalance returns balance of Wei (ether) in the account
func (u User) GetWeiBalance() *big.Int {
	balance, err := u.client.BalanceAt(context.Background(), u.address, nil)
	utils.CheckError(err)

	return balance
}

// GetDaiBalance returns balance of dai stable coin in the account
func (u User) GetDaiBalance() *big.Int {
	daiAddress := common.HexToAddress(daiContract)
	return u.getTokenBalance(daiAddress)
}

// GetMkrBalance returns the balance of mkr in the account
func (u User) GetMkrBalance() *big.Int {
	mkrAddress := common.HexToAddress(mkrContract)
	return u.getTokenBalance(mkrAddress)
}

// GetSaiBalance returns the balance of sai in the account
func (u User) GetSaiBalance() *big.Int {
	sntAddress := common.HexToAddress(saiContract)
	return u.getTokenBalance(sntAddress)
}

func (u User) getTokenBalance(tokenAddress common.Address) *big.Int {
	instance, err := token.NewToken(tokenAddress, u.client)
	utils.CheckError(err)

	bal, err := instance.BalanceOf(&bind.CallOpts{}, u.address)
	utils.CheckError(err)

	return bal
}

type valueFunc func(user User) *big.Int

type coinValue struct {
	coin    string
	balance *big.Int
}

// GetBalances returns a map with the balances of sai, mkr, wei and dai in the account
func (u User) GetBalances() map[string]*big.Int {
	tokens := map[string]valueFunc{
		"sai": func(user User) *big.Int { return user.GetSaiBalance() },
		"mkr": func(user User) *big.Int { return user.GetMkrBalance() },
		"wei": func(user User) *big.Int { return user.GetWeiBalance() },
		"dai": func(user User) *big.Int { return user.GetDaiBalance() },
	}

	c := make(chan coinValue)

	for coin, invocation := range tokens {
		go u.fetchBalance(coin, invocation, c)
	}

	balances := make(map[string]*big.Int)

	for range tokens {
		balance := <-c
		balances[balance.coin] = balance.balance
	}

	return balances
}

func (u User) fetchBalance(coin string, invocation valueFunc, c chan coinValue) {
	balance := invocation(u)
	c <- coinValue{coin: coin, balance: balance}
}
