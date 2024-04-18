package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
)

func main() {
	// MainprivateKey := os.Getenv("PRIVATE_KEY")
	MainprivateKey := "7cf14cea97ee9ace52e647f282686ef9bd32c5fd686272190520289905885293"
	privateKey, err:= crypto.GenerateKey()
	// publicKey:= privateKey.Public()
	if err != nil {
		log.Fatal(err.Error())
	}

	rawPrivateKey:= crypto.FromECDSA(privateKey)
	mainInBytes:= common.Hex2Bytes(MainprivateKey)
	ZkSyncEraProvider := "https://sepolia.era.zksync.dev"
	EthereumProvider := "wss://eth-sepolia.g.alchemy.com/v2/VtFb4uQ7Vc5l414EGwXDDzcSClkHv9TY"

	zkClient, err := clients.Dial(ZkSyncEraProvider)
	if err != nil {
		log.Fatal(err.Error())
	}

	ethClient, err:= ethclient.Dial(EthereumProvider)
	if err != nil {
		log.Fatal(err.Error())
	}

	ethchainID, err:= ethClient.ChainID(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("ETH Chian ID:", ethchainID)

	zkchaiID, err:= zkClient.ChainID(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("ZK Chain ID:", zkchaiID)

	zkethchainID,err:= zkClient.L1ChainID(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("zkEth Chain ID:", zkethchainID)

	wallet, err:= accounts.NewWallet(mainInBytes, &zkClient, ethClient)
	if err != nil {
		log.Fatal(err.Error())
	}

	balance, err:= wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Balance of the address is:", balance)

	fmt.Println("New Private Key:", privateKey)
	fmt.Println("Old private key: ", MainprivateKey)
	fmt.Println("Raw Private Key:", rawPrivateKey)
	fmt.Println("Main in Bytes", mainInBytes)
}