package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
)

func main() {
	prune := flag.Bool("prune", false, "drop database before running app")
	devMode := flag.Bool("dev", false, "disable signature verification")
	flag.Parse()

	var cfg *config.Config
	if *devMode {
		cfg = config.GetTestConfig()
	} else {
		cfg = config.GetConfig()
	}

	cfg.Rollup.Prune = *prune

	cmd := commander.NewCommander(cfg)

	setupCloseHandler(cmd)

	err := cmd.StartAndWait()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func setupCloseHandler(cmd *commander.Commander) {
	c := make(chan os.Signal)
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
