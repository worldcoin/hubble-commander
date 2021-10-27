package main

import (
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/utils"
	log "github.com/sirupsen/logrus"
)

func deployCommanderContracts(filePath string) error {
	cfg := config.GetDeployerConfigAndSetupLogger()
	blockchain, err := commander.GetChainConnection(cfg.Ethereum)
	if err != nil {
		return err
	}

	chainSpec, err := commander.Deploy(cfg, blockchain)
	if err != nil {
		return err
	}
	log.Printf(*chainSpec)

	return utils.StoreChainSpec(filePath, *chainSpec)
}
