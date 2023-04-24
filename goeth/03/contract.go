package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func callContract(client *ethclient.Client) {
	addr := deployContracts(client)
	searchContract(addr, client)
}

func deployContracts(client *ethclient.Client) (contractAddress common.Address) {
	// 1. get nonce
	// 2. get gas price
	// 3. get gas limit
	// 4. get contract data
	// 5. get chain id
	// 6. sign
	// 7. send
	// 8. get tx hash
	// 9. get receipt
	// 10. get contract address
	privateKey, err := crypto.HexToECDSA("3c6af19adf147b5b2f6111c4015996016f2b3f92aaf0457bff6d053945607e6c")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("fromAddress", fromAddress.Hex())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("nonce", nonce)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	addr, _, _, err := DeployStore(auth, client, "1.0")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("contract address", addr.Hex())
	contractAddress = addr

	return
}

func searchContract(addr common.Address, client *ethclient.Client) (contract *Store) {
	contract, err := NewStore(addr, client)
	if err != nil {
		log.Fatal(err)
	}

	version, err := contract.Version(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("contract version: %v\n", version)

	return contract
}
