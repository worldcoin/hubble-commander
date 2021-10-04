package main

import (
	"github.com/Worldcoin/hubble-commander/scripts"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := scripts.QueryRegisteredPublicKeys()
	if err != nil {
		log.Fatal(err)
	}
}
