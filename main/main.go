package main

import (
	"log"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
)

func main() {
	cfg := config.GetConfig()

	go commander.RollupLoop(&cfg)

	log.Fatal(api.StartAPIServer(&cfg))
}
