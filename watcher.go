package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Watcher struct {
	Client *ethclient.Client
	Config *EventConfig
}

func newWatcher(client *ethclient.Client, config *EventConfig) *Watcher {
	return &Watcher{
		Client: client,
		Config: config,
	}
}

func (w *Watcher) Run() {
	addr := common.HexToAddress(w.Config.ContractAddress)
	eventTopicId := w.Config.ABI.Events[w.Config.EventName].ID
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(w.Config.StartingBlock),
		Addresses: []common.Address{addr},
		Topics:    [][]common.Hash{[]common.Hash{eventTopicId}},
	}
	fmt.Printf("TOPIC: %+v\n", [][]common.Hash{[]common.Hash{eventTopicId}})

	logs, err := w.Client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("--- logs count %+v\n", len(logs))

	for _, vLog := range logs {
		ptr, err := parseEvent(w.Config, vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ptr: %+v\n", ptr)
	}
}
