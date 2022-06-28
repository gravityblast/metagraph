package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/ethclient"
)

// FIXME: create a personal dedicated on pinata.cloud
const IPFS_GATEWAY_URL = "https://ipfs.io/ipfs"
const IPFS_CLIENT_TIMEOUT = 10
const PIN_WORKERS = 6

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

	ethClient, err := ethclient.Dial(config.ProviderURL)
	if err != nil {
		log.Fatal(err)
	}

	ptrQueue := make(chan *MetaPointer, 10)
	quit := make(chan struct{})

	for i := 0; i < PIN_WORKERS; i++ {
		go pinWorker(i, ptrQueue)
	}

	for _, ec := range config.Events {
		w := newWatcher(ethClient, ec, ptrQueue)
		w.Run()
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		quit <- struct{}{}
	}()

	<-quit
}
