package main

import (
	"flag"
	"log"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
)

func main() {
	prune := flag.Bool("prune", false, "drop database before running app")
	devMode := flag.Bool("dev", false, "disable signature verification")
	flag.Parse()

	var cfg config.Config
	if *devMode {
		cfg = config.GetTestConfig()
	} else {
		cfg = config.GetConfig()
	}

	cfg.Rollup.Prune = *prune

	cmd := commander.NewCommander(&cfg)

	err := cmd.Start()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
