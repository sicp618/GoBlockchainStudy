package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func callListenEvent() {
	host := "ws://localhost:7545"
	client, err := ethclient.Dial(host)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			fmt.Println("sub err", err)
		case vLog := <-logs:
			fmt.Println("log", vLog)
		}
	}
}

func readEvents(client *ethclient.Client) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(1000),
		Addresses: []common.Address{contractAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(StoreABI)))
	if err != nil {
		log.Fatal(err)
	}

	type Event struct {
		Key   [32]byte
		Value [32]byte
	}

	for _, vLog := range logs {
		fmt.Println("read log", vLog.Data)
		var event Event
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("event", string(event.Key[:]), string(event.Value[:]))
	}
}
