package main

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/scripts"
	"github.com/urfave/cli/v2"
)

var exportTypes = []string{"state", "accounts"}

func exportData(ctx *cli.Context) error {
	file := ctx.String("file")

	var err error
	switch ctx.String("type") {
	case exportTypes[0]:
		err = scripts.ExportStateLeaves(file)
	case exportTypes[1]:
		err = scripts.ExportAccounts(file)
	default:
		return fmt.Errorf("invalid export data type, supported: %v", exportTypes)
	}
	return err
}
