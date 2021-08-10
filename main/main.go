package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := getConfig()

	configureLogger(cfg)
	logConfig(cfg)
	cmd := commander.NewCommander(cfg)

	setupCloseHandler(cmd)

	err := cmd.StartAndWait()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func getConfig() *config.Config {
	prune := flag.Bool("prune", false, "drop database before running app")
	disableSignatures := flag.Bool("disable-signatures", false, "disable signature verification")
	flag.Parse()

	cfg := config.GetConfig()
	cfg.Bootstrap.Prune = *prune
	cfg.Rollup.DisableSignatures = *disableSignatures
	return cfg
}

func configureLogger(cfg *config.Config) {
	if cfg.Log.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(cfg.Log.Level)
}

func logConfig(cfg *config.Config) {
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	log.Debugf("Loaded config: %s", string(jsonCfg))
}

func setupCloseHandler(cmd *commander.Commander) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Warning("Stopping commander gracefully...")
		err := cmd.Stop()
		if err != nil {
			log.Errorf("Error while stopping: %+v", err)
		}
	}()
}
