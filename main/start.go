package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	log "github.com/sirupsen/logrus"
)

func startCommander() error {
	cfg := config.GetCommanderConfigAndSetupLogger()
	blockchain, err := commander.GetChainConnection(cfg.Ethereum)
	if err != nil {
		return err
	}

	cmd := commander.NewCommander(cfg, blockchain)
	setupCloseHandler(cmd)

	return cmd.StartAndWait()
}

func setupCloseHandler(cmd *commander.Commander) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Warning("Stopping commander gracefully...")
		err := cmd.Stop()
		if err != nil {
			log.Panicf("Failed to stop commander gracefully: %+v", err)
		}
	}()
}
