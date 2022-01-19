package main

import (
	"flag"

	"github.com/Worldcoin/hubble-commander/scripts"
	log "github.com/sirupsen/logrus"
)

var (
	deployFlagSet = flag.NewFlagSet(deploy, flag.ExitOnError)
	chainSpecFile = deployFlagSet.String("file", "chain-spec.yaml", "target file to save the chain spec to")

	exportFlagSet = flag.NewFlagSet(export, flag.ExitOnError)
	exportType    = exportFlagSet.String("type", "state", "type of data to export")
	exportFile    = exportFlagSet.String("file", "exported-data.json", "target file to save exported data to")
)

func handleStartCommand() {
	err := startCommander()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func handleDeployCommand(args []string) {
	err := deployFlagSet.Parse(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = deployCommanderContracts(*chainSpecFile)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func handleExportCommand(args []string) {
	err := exportFlagSet.Parse(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	switch *exportType {
	case "state":
		err = scripts.ExportStateLeaves(*exportFile)
	default:
		exitWithHelpMessage()
	}
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
