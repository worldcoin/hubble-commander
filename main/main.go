package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	args := os.Args

	// TODO-CHAIN handle flags with Go's FlagSets
	if len(args) == 1 {
		log.Fatal("please provide an arg")
	}

	if args[1] == "deploy" {
		if len(args) == 3 {
			deployCommanderContracts(args[2])
		} else {
			log.Fatal("please provide a filename")
		}
	}

	if args[1] == "start" {
		startCommander()
	}
}

func setupCommander() *commander.Commander {
	cfg := getConfig()

	configureLogger(cfg)
	logConfig(cfg)
	cmd := commander.NewCommander(cfg)

	setupCloseHandler(cmd)

	return cmd
}

func startCommander() {
	cmd := setupCommander()
	err := cmd.StartAndWait()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func deployCommanderContracts(filename string) {
	cmd := setupCommander()
	chainSpec, err := cmd.Deploy()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = ioutil.WriteFile(filename, []byte(*chainSpec), 0644)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	log.Printf(*chainSpec)
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
