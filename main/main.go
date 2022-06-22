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
