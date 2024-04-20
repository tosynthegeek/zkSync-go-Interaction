package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
)

func main() {
	privateKey:= os.Getenv("PRIVATE_KEY")
	zkSyncEraProvider := "https://sepolia.era.zksync.dev"
	ethProvider:= ""

	// Connect to the zkSync network 
	ZKClient, err := clients.Dial(zkSyncEraProvider)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ZKClient.Close() // Close on exit

	// Connect to the Ethereum network
	ethClient, err:= ethclient.Dial(ethProvider)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ethClient.Close() // Close on exit

	fmt.Println("Clients created......")
	
	// Create new wallet using private key and network clients
	wallet, err:= accounts.NewWallet(common.Hex2Bytes(privateKey), &ZKClient, ethClient)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Wallet Created.... ")

	// Get balance before deposit. Setting block number to nil so it reurns balance at the latest block
	balance, err:= wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Balance before deposit is: ", balance)

	// Initiate withdrawal from zkSync(L2) to Ethereum(L1)
	tx, err:= wallet.Withdraw(nil, accounts.WithdrawalTransaction{
		Amount: big.NewInt(1000000000),
		Token: utils.EthAddress,
		To: wallet.Address(),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Withdrawal Hash: ", tx.Hash())

	// Get balance after deposit. Setting block number to nil so it reurns balance at the latest block
	balance, err= wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Balance after deposit:", balance)

}
