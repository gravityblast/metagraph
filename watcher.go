package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Watcher struct {
	Client *ethclient.Client
	Config *EventConfig
	Logger *logrus.Entry
}

func newWatcher(client *ethclient.Client, config *EventConfig) *Watcher {
	return &Watcher{
		Client: client,
		Config: config,
		Logger: logger.WithField("object", "watcher"),
	}
}

func (w *Watcher) Run() {
	runUUID := uuid.New()
	l := w.Logger.WithFields(logrus.Fields{
		"action":  "run",
		"runUUID": runUUID.String(),
	})
	addr := common.HexToAddress(w.Config.ContractAddress)
	eventTopicId := w.Config.ABI.Events[w.Config.EventName].ID
	fromBlock := big.NewInt(w.Config.StartingBlock)
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: []common.Address{addr},
		Topics:    [][]common.Hash{[]common.Hash{eventTopicId}},
	}

	l.WithFields(logrus.Fields{
		"fromBlock": fromBlock.String(),
		"address":   addr,
		"topic":     eventTopicId,
	}).Info("FilterLogs")

	logs, err := w.Client.FilterLogs(context.Background(), query)
	if err != nil {
		l.Error(err.Error())
		return
	}

	l.WithFields(logrus.Fields{
		"resultsCount": len(logs),
	}).Info("FilterLogs results")

	for _, vLog := range logs {
		ptr, err := parseEvent(w.Config, vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		l.WithFields(logrus.Fields{
			"metadataProtocol": ptr.Protocol,
			"metadataPointer":  ptr.Pointer,
		}).Info("metaPointer found")
	}
}
