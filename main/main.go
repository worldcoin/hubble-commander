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
	devMode := flag.Bool("dev", false, "disable signature verification")
	flag.Parse()

	var cfg *config.Config
	if *devMode {
		cfg = config.GetTestConfig()
	} else {
		cfg = config.GetConfig()
	}
	cfg.Bootstrap.Prune = *prune
	return cfg
}

func configureLogger(cfg *config.Config) {
	if cfg.Log.Format == "json" {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
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
		log.Println("\nStopping commander gracefully...")
		err := cmd.Stop()
		if err != nil {
			log.Printf("Error while stopping: %+v", err)
		}
	}()
}
