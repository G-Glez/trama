package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"trama/cmd/lambda/provider"
)

type Handler struct {
	adapter *ginadapter.GinLambdaV2
}

func New(p *provider.Provider) *Handler {
	return &Handler{adapter: ginadapter.NewV2(p.Gin())}
}

func (h *Handler) HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return h.adapter.ProxyWithContext(ctx, req)
}
