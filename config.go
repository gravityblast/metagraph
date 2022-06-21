package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type EventConfig struct {
	ID               string `json:"id"`
	Description      string `json:"description"`
	ContractName     string `json:"contractName"`
	ContractAddress  string `json:"contractAddress"`
	EventName        string `json:"eventName"`
	MetaPointerField string `json:"metaPointerField"`
	ProtocolField    string `json:"protocolField"`
	PointerField     string `json:"pointerField"`
	StartingBlock    int64  `json:"startingBlock"`

	ABI abi.ABI
}

type Config struct {
	ProviderURL string         `json:"providerURL"`
	Events      []*EventConfig `json:"events"`
}

func parseConfig(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, err
	}

	for _, event := range config.Events {
		abiContent, err := os.ReadFile(filepath.Join("abis", fmt.Sprintf("%s.json", event.ContractName)))
		if err != nil {
			return nil, err
		}

		abi, err := abi.JSON(bytes.NewReader(abiContent))
		if err != nil {
			return nil, err
		}

		event.ABI = abi
	}

	return config, err
}
