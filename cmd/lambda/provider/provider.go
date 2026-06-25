package provider

import (
	"log/slog"
	"os"
	"trama/cmd/lambda/handler"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/caarlos0/env/v11"
)

type envConfig struct {
	Env                    string `env:"ENV,required"`
	JWTSecret              string `env:"JWT_SECRET,required"`
	DynamoDBUsersTableName string `env:"DYNAMODB_USERS_TABLE_NAME,required"`
}

type Provider struct {
	env envConfig
	InfraProvider
	ApiProvider
}

func NewProvisionedProvider() *Provider {
	var e envConfig
	if err := env.Parse(&e); err != nil {
		panic(err)
	}

	p := &Provider{env: e}
	p.provisionLogger()
	p.provisionInfra()
	p.provisionApi()

	return p
}

func (p *Provider) Handler() *handler.Handler {
	p.router.Setup(p.gin, p.authMiddleware)
	adapter := ginadapter.NewV2(p.gin)
	return handler.New(adapter)
}


func (p *Provider) provisionLogger() {
	level := slog.LevelInfo
	if p.env.Env == "dev" {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}
