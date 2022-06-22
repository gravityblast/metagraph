package main

import (
	"log"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/ethclient"
)

var logger *logrus.Entry

func init() {
	logger = logrus.New().WithField("app", "metagraph")
	// log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	config, err := parseConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial(config.ProviderURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, ec := range config.Events {
		w := newWatcher(client, ec)
		w.Run()
	}
}
