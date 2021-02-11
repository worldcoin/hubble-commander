package main

import (
	"log"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/config"
)

func main() {
	cfg := config.GetConfig()
	log.Fatal(api.StartApiServer(&cfg))
}
