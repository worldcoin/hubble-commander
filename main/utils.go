package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/Worldcoin/hubble-commander/bls"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func newWallet(ctx *cli.Context) error {
	privateKey := make([]byte, 32)
	_, err := rand.Read(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := bls.NewWallet(privateKey, bls.Domain{0x00, 0x00, 0x00, 0x00})
	if err != nil {
		log.Fatal(err)
	}

	result, _ := json.Marshal(
		struct {
			PrivateKey string
			PublicKey  string
		}{
			PrivateKey: fmt.Sprintf("0x%x", privateKey),
			PublicKey:  wallet.PublicKey().String(),
		},
	)

	fmt.Printf("%s\n", result)
	return nil
}
