package main

import (
	"log"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/config"
)

func main() {
	cfg := config.CreateConfig(
		"dev-0.1.0",
		8080,
		"hubble_test",
		"hubble",
		"root",
	)
	log.Fatal(api.StartApiServer(cfg))
}
