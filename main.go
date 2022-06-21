package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	config, err := parseConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial(config.ProviderURL)

	for _, ec := range config.Events {
		w := newWatcher(client, ec)
		fmt.Printf("-------------------- %+v\n", ec.Description)
		w.Run()
	}
}
