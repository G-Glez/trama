package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"trama/cmd/lambda/provider"
)

func main() {
	p := provider.NewProvisionedProvider()
	h := p.Handler()
	lambda.Start(h.HandleRequest)
}
