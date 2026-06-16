package main

import (
	"os"

	"trama/cmd/cli/provider"
)

var (
	p = provider.NewProvisionedProvider()
)

func main() {
	if err := p.CLI().Start(); err != nil {
		os.Exit(1)
	}
}
