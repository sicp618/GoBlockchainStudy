package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	host := "https://cloudflare-eth.com"
	host = "http://localhost:7545"
	client, err := ethclient.Dial(host)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections
}
