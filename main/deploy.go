package main

import (
	"os"
	"path/filepath"

	"github.com/Worldcoin/hubble-commander/commander"
	log "github.com/sirupsen/logrus"
)

func deployCommanderContracts(filePath string) error {
	cfg := getConfigAndSetupLogger()
	chain, err := commander.GetChainConnection(cfg.Ethereum)
	if err != nil {
		return err
	}

	chainSpec, err := commander.Deploy(cfg, chain)
	if err != nil {
		return err
	}
	log.Printf(*chainSpec)

	return storeChainSpec(filePath, *chainSpec)
}

func storeChainSpec(filePath, chainSpec string) error {
	dirPath := filepath.Dir(filePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(chainSpec), 0600)
}
