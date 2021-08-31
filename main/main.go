package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	deployCommand = flag.NewFlagSet("deploy", flag.ExitOnError)
	chainSpecFile = deployCommand.String("file", "chain-spec.yaml", "target file to save the chain spec to")
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
	case "start":
		handleStartCommand()
	case "deploy":
		handleDeployCommand(os.Args[2:])
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
