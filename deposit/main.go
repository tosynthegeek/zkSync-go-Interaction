package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
)

func main() {
	privateKey:= os.Getenv("PRIVATE_KEY")
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

	// Get balance before deposit. Setting block number to nil so it reurns balance at the latest block
	balance, err:= wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Balance before deposit is: ", balance)

	// Deposit token from L1 to L2
	tx, err:= wallet.Deposit(nil, accounts.DepositTransaction{
		Token: utils.EthAddress,
		Amount: big.NewInt(1000000000), // In Wei
		To: wallet.Address(), // Deposits to our own address
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("L1 Tx Hash: ", tx.Hash())

	// Wait for deposit transaction to be finalized on L1 (Ethereum) network
	fmt.Println("Waiting for tx to be finished on L1...")
	_, err = bind.WaitMined(context.Background(), ethClient, tx)
	if err != nil {
		log.Panic(err)
	}
	
	// Get reciept for Ethereum Transaction
	ethReciept, err:= ethClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Reciept of Ethereum Tx: ", ethReciept)

	// Use reiept from Ethereum(L1) to return transaction on zkSync(L2)
	zkTx, err:= ZKClient.L2TransactionFromPriorityOp(context.Background(), ethReciept)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("L2 Hash:", zkTx)

	// Wait for deposit transaction to be finalized on L2 (zkSync) network (can take 5-7 minutes)
	_, err = ZKClient.WaitMined(context.Background(), zkTx.Hash)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get balance after deposit. Setting block number to nil so it reurns balance at the latest block
	balance, err= wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Balance after deposit:", balance)
}