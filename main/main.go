package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	deployCommand = flag.NewFlagSet("deploy", flag.ExitOnError)
	chainSpecFile = deployCommand.String("file", "chain-spec.yaml", "target file to save the chain spec to")
)

func exitWithHelpMessage() {
	fmt.Printf("Subcommand required:\n" +
		"start - starts the commander\n" +
		"deploy - deploys contracts and saves chain spec. Usage:\n")
	deployCommand.PrintDefaults()
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		exitWithHelpMessage()
	}

	switch os.Args[1] {
	case "start":
		startCommander()
	case "deploy":
		handleDeployCommand(os.Args[2:])
	default:
		exitWithHelpMessage()
	}
}

func handleDeployCommand(args []string) {
	err := deployCommand.Parse(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	deployCommanderContracts(*chainSpecFile)
}

func setupCommander() *commander.Commander {
	cfg := config.GetConfig()

	configureLogger(cfg)
	logConfig(cfg)

	chain, err := commander.GetChainConnection(cfg.Ethereum)
	if err != nil {
		log.Fatal(err)
	}
	cmd := commander.NewCommander(cfg, chain)

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
	dirName := filepath.Dir(filename)
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = os.WriteFile(filename, []byte(*chainSpec), 0600)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	log.Printf(*chainSpec)
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
