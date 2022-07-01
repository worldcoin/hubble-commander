package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "hubble-commander",
		Usage: "Commander for Hubble optimistic rollup",
		Commands: cli.Commands{
			{
				Name:   "start",
				Usage:  "start the commander",
				Action: startCommander,
			},
			{
				Name:   "auditDatabase",
				Usage:  "which prefixes and consuming the most space?",
				Action: auditDatabase,
			},
			{
				Name:   "newWallet",
				Usage:  "create a new BLS wallet",
				Action: newWallet,
			},
			{
				Name:   "sendTransaction",
				Usage:  "send tokens from one hubble account to another",
				Action: sendTransaction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "rpcurl",
						Usage: "location of the hubble commander",
						Value: "http://localhost:8080",
					},
					&cli.StringFlag{
						Name:  "type",
						Usage: "either TRANSFER or CREATE2TRANSFER",
						Value: "TRANSFER",
					},
					&cli.StringFlag{
						Name:     "privateKey",
						Aliases:  []string{"privatekey"},
						Usage:    "The hex-encoded private key",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "from",
						Usage:    "which wallet is sending",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to",
						Usage:    "which wallet is receiving",
						Required: true,
					},
					&cli.Uint64Flag{
						Name:     "amount",
						Usage:    "how much to send",
						Required: true,
					},
					&cli.Uint64Flag{
						Name:     "fee",
						Usage:    "how much to pay the sequencer",
						Required: true,
					},
				},
			},
			{
				Name:   "benchmark",
				Usage:  "run transactions against a commander",
				Action: benchmarkHubble,
			},
			{
				Name:  "deploy",
				Usage: "deploy contracts and save chain spec",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "file",
						Usage: "target file to save the chain spec to",
						Value: "chain-spec.yaml",
					},
				},
				Action: deployContracts,
			},
			{
				Name:  "export",
				Usage: "export data to file in json format",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "type",
						Usage:    "type of data to export",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "file",
						Usage: "target file to save exported data to",
						Value: "exported-data.json",
					},
				},
				Action: exportData,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
