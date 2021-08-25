package main

import (
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/utils"
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

	return utils.StoreChainSpec(filePath, *chainSpec)
}
