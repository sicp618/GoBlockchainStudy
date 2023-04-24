package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)


func main() {
	host := "http://localhost:7545"
	// host = "ws://localhost:8545"
	// addr := "0xEf41Deb6F24d50887Ed8c15a24862A7d5029Cf66"
	client, err := ethclient.Dial(host)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// readTxInfo(client)
	// listenNewTx(client)
	callContract(client)
}

func listenNewTx(client *ethclient.Client) {
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal("sub err", err)
	}
	fmt.Println("begin listen")
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("select err")
			fmt.Println("sub err", err)
		case header := <-headers:
			fmt.Println("select header")
			hash := header.Hash()
			fmt.Println(hash)
			block, err := client.BlockByHash(context.Background(), hash)
			if err != nil {
				fmt.Println("find block err", err)
				break
			}

			fmt.Println("hex", block.Hash().Hex())
			fmt.Println("number", block.Number().Uint64())
			fmt.Println("time", block.Time())
			fmt.Println("nonce", block.Nonce())
			fmt.Println("tx length", len(block.Transactions()))
		}
	}
}

func readTxInfo(client *ethclient.Client) {
	header, _ := client.HeaderByNumber(context.Background(), nil)
	fmt.Println(header.Number.String())

	number := header.Number
	// number = big.NewInt(143)
	block, _ := client.BlockByNumber(context.Background(), number)
	tc, _ := client.TransactionCount(context.Background(), block.Hash())
	fmt.Println(tc)

	chainID, _ := client.NetworkID(context.Background())
	fmt.Println("chainID", chainID)
	for _, tx := range block.Transactions() {
		fmt.Println("TX", tx.Hash().Hex())
		fmt.Println("GAS", tx.Gas())
		fmt.Println("ChainID", tx.ChainId().String())
		fmt.Println("nonce", tx.Nonce())
		fmt.Println("to", tx.To())
		fmt.Println("data", tx.Data())

		if address, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx); err == nil {
			fmt.Println(address.Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		} else {
			fmt.Println(err)
		}
	}
}
