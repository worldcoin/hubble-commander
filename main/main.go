package main

import (
	"log"

	"github.com/Worldcoin/hubble-commander/api"
)

func main() {
	log.Fatal(api.StartApiServer(":8080"))
}
