package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Worldcoin/hubble-commander/scripts"
	log "github.com/sirupsen/logrus"
)

const (
	start  = "start"
	deploy = "deploy"
	export = "export"
)

var (
	deployCommand = flag.NewFlagSet(deploy, flag.ExitOnError)
	chainSpecFile = deployCommand.String("file", "chain-spec.yaml", "target file to save the chain spec to")

	exportCommand = flag.NewFlagSet(export, flag.ExitOnError)
	exportType    = exportCommand.String("type", "state", "type of data to export")
	exportFile    = exportCommand.String("file", "exported-data.json", "target file to save exported data to")
)

func exitWithHelpMessage() {
	fmt.Printf("Subcommand required:\n" +
		"start - starts the commander\n" +
		"deploy - deploys contracts and saves chain spec. Usage:\n")
	deployCommand.PrintDefaults()
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		exitWithHelpMessage()
	}

	switch os.Args[1] {
	case start:
		handleStartCommand()
	case deploy:
		handleDeployCommand(os.Args[2:])
	case export:
		handleExportCommand(os.Args[2:])
	default:
		exitWithHelpMessage()
	}
}

func handleStartCommand() {
	err := startCommander()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func handleDeployCommand(args []string) {
	err := deployCommand.Parse(args)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = deployCommanderContracts(*chainSpecFile)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func handleExportCommand(args []string) {
	err := exportCommand.Parse(args)
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
