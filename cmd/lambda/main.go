package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"trama/cmd/lambda/handler"
	"trama/cmd/lambda/provider"
)

var (
	p       = provider.NewProvisionedProvider()
	h       = handler.New(p)
)

func main() {
	defer p.Close()
	lambda.Start(h.HandleRequest)
}
