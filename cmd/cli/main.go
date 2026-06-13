package main

import (
	"os"

	"trama/cmd/cli/provider"
)

var (
	p     = provider.NewProvider()
	dbCon = p.DBCon()
	cli   = p.Cli(dbCon)
)

func main() {
	if err := cli.Start(); err != nil {
		os.Exit(1)
	}
}
