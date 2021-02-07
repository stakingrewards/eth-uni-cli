package wallet

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/Bullrich/eth-uni-cli/uniswap"
	"github.com/Bullrich/eth-uni-cli/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"time"
)

const uniswapRouterAddress = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"

func (u *User) obtainContractsRequirements(privateWalletKey string, weiToUse int64) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(privateWalletKey)
	utils.CheckError(err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := u.client.PendingNonceAt(context.Background(), fromAddress)
	utils.CheckError(err)

	gasPrice, err := u.client.SuggestGasPrice(context.Background())
	utils.CheckError(err)

	auth := bind.NewKeyedTransactor(privateKey)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(weiToUse) // in wei
	auth.GasLimit = uint64(300000)    // in units
	auth.GasPrice = gasPrice

	return auth
}

func BuyCoin(user *User, privateKey string, coinName string, input int64, output int64) (string, error) {
	coinAddress, err := getCoinAddress(coinName)
	if err != nil {
		return "", err
	}

	auth := user.obtainContractsRequirements(privateKey, input)

	coinContract := common.HexToAddress(coinAddress)

	transactionHex := user.buyCoin(coinContract, output, auth)
	return transactionHex, err
}

func getCoinAddress(coinName string) (string, error) {
	switch coinName {
	case "dai":
		return daiContract, nil
	case "mkr":
		return mkrContract, nil
	case "sai":
		return saiContract, nil
	}

	return "", errors.New(fmt.Sprintf("Coin %v not valid. Use 'dai', 'mkr' or 'sai'", coinName))
}

func BuyDai(user *User, privateKey string) {
	auth := user.obtainContractsRequirements(privateKey, 10000000000000000)

	address := common.HexToAddress(daiContract)

	user.buyCoin(address, 1560000, auth)
}

func (u *User) buyCoin(coinAddress common.Address, amountToBuy int64, auth *bind.TransactOpts) string {
	contractAddress := common.HexToAddress(uniswapRouterAddress)
	uniswapContract, err := uniswap.NewUniswap(contractAddress, u.client)
	utils.CheckError(err)

	amountOut := big.NewInt(amountToBuy)
	callOptsParams := &bind.CallOpts{Pending: true}
	weth, err := uniswapContract.WETH(callOptsParams)
	utils.CheckError(err)
	addresses := []common.Address{
		weth,
		coinAddress,
	}

	now := time.Now().Unix()
	i := new(big.Int).SetInt64(now + (time.Second.Milliseconds() * 40))

	tx, err := uniswapContract.SwapETHForExactTokens(auth, amountOut, addresses, u.address, i)
	utils.CheckError(err)

	fmt.Printf("tx sent. See the transaction in: https://etherscan.io/tx/%s\n", tx.Hash().Hex())

	return tx.Hash().Hex()
}
