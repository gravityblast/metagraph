package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// FIXME: create a personal dedicated on pinata.cloud
const IPFS_GATEWAY_URL = "https://ipfs.io/ipfs"
const IPFS_CLIENT_TIMEOUT = 5

type Watcher struct {
	Client     *ethclient.Client
	Config     *EventConfig
	Logger     *logrus.Entry
	ipfsClient *http.Client
}

func newWatcher(client *ethclient.Client, config *EventConfig) *Watcher {
	ipfsClient := &http.Client{
		Timeout: IPFS_CLIENT_TIMEOUT * time.Second,
	}

	return &Watcher{
		Client:     client,
		Config:     config,
		Logger:     logger.WithField("object", "watcher"),
		ipfsClient: ipfsClient,
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
			l.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("error parsing event")
		}

		l.WithFields(logrus.Fields{
			"metadataProtocol": ptr.Protocol,
			"metadataPointer":  ptr.Pointer,
		}).Info("metaPointer found")

		if len(w.Config.Embedded) > 0 {
			w.parseEmbedded(runUUID, ptr)
		}
	}
}

func (w *Watcher) parseEmbedded(runUUID uuid.UUID, ptr *MetaPointer) {
	l := w.Logger.WithFields(logrus.Fields{
		"action":           "parseEmbedded",
		"runUUID":          runUUID.String(),
		"metadataProtocol": ptr.Protocol,
		"metadataPointer":  ptr.Pointer,
	})

	l.Info("fetching")
	u := fmt.Sprintf("%s/%s", IPFS_GATEWAY_URL, ptr.Pointer)

	resp, err := w.ipfsClient.Get(u)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error fetching")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.WithFields(logrus.Fields{
			"status":     resp.Status,
			"statusCode": resp.StatusCode,
		}).Error("unexpected status code")
		return
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error reading content")
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error parsing json")
		return
	}

	for _, embed := range w.Config.Embedded {
		ptr, err := parseMetaPointerFromMap(embed, data)
		if err != nil {
			l.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("error parsing event")
			continue
		}

		l.WithFields(logrus.Fields{
			"metadataProtocol": ptr.Protocol,
			"metadataPointer":  ptr.Pointer,
		}).Info("metaPointer found")
	}
}
