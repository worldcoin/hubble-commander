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
	if len(os.Args) < 2 {
		log.Fatal("start or deploy subcommand is required")
	}

	switch os.Args[1] {
	case "start":
		handleStartCommand(os.Args[2:])
	case "deploy":
		handleDeployCommand(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func handleStartCommand(args []string) {
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	err := startCommand.Parse(args)
	if err != nil {
		log.Fatal(err)
	}
	startCommander()
}

func handleDeployCommand(args []string) {
	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	chainSpecFile := deployCommand.String("file", "chain-spec.yaml", "TODO")
	err := deployCommand.Parse(args)
	if err != nil {
		log.Fatal(err)
	}
	deployCommanderContracts(*chainSpecFile)
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
	return config.GetConfig()
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
