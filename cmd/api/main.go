package main

import (
	"log"

	"trama/cmd/api/provider"
)

var (
	p = provider.NewProvisionedProvider()
)

func main() {
	defer p.Close()

	log.Printf("Starting server on port %s", p.Config().Port)
	p.Router().Run(":" + p.Config().Port)
}
