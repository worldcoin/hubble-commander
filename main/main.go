package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := getConfig()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

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

func setupCloseHandler(cmd *commander.Commander) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nStopping commander gracefully...")
		err := cmd.Stop()
		if err != nil {
			fmt.Printf("Error while stopping: %+v", err)
		}
	}()
}
