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
	toAddress:= common.HexToAddress("0xD109E8C395741b4b3130E3D84041F8F62aF765Ef")
	zkSyncEraProvider := "https://sepolia.era.zksync.dev"
	ethProvider:= "wss://eth-sepolia.g.alchemy.com/v2/VtFb4uQ7Vc5l414EGwXDDzcSClkHv9TY"

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

	// Get wallet balance before transfer
	myBalance, err:= wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("My balance before transfer is: ", myBalance)

	// Check recipient balance before transfer
	recipientBalance, err:= ZKClient.BalanceAt(context.Background(), toAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Recipient balance before transfer is: ", recipientBalance)

	// Transfer token from my wallet to recipient address
	tx, err:= wallet.Transfer(nil, accounts.TransferTransaction{
		Token: utils.EthAddress,
		Amount: big.NewInt(1000000000),
		To: toAddress,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Transaction Hash: ", tx.Hash())

	_, err = ZKClient.WaitMined(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err.Error())
	}

	// Check wallet balance after transfer
	myBalance, err = wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("My balance after transfer is: ", myBalance)

	// Check recipient balance after transfer
	recipientBalance, err = ZKClient.BalanceAt(context.Background(), toAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Recipient balance after transfer is: ", recipientBalance)

}