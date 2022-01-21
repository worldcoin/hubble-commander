package main

import (
	"fmt"
	"os"
)

const (
	start  = "start"
	deploy = "deploy"
	export = "export"
)

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

func exitWithHelpMessage() {
	fmt.Printf(`
Available subcomands:
start - starts the commander
deploy - deploys contracts and saves chain spec		
`)
	deployFlagSet.Usage()

	fmt.Println("export - exports data to file in json format")
	exportFlagSet.Usage()
	os.Exit(1)
}
