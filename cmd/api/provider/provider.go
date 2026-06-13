package provider

import (
	"database/sql"
	"log"

	"github.com/caarlos0/env/v11"

	"trama/internal/api/config"
	"trama/pkg/dbcon"
)

type envConfig struct {
	Port         string `env:"PORT,required"`
	GinMode      string `env:"GIN_MODE,required"`
	DatabasePath string `env:"DATABASE_PATH,required"`
}

type Provider struct {
	env envConfig
}

func NewProvider() *Provider {
	var e envConfig
	if err := env.Parse(&e); err != nil {
		panic(err)
	}

	return &Provider{env: e}
}

func (p *Provider) Config() *config.Config {
	return &config.Config{
		Port:         p.env.Port,
		GinMode:      p.env.GinMode,
		DatabasePath: p.env.DatabasePath,
	}
}

func (p *Provider) DB() *sql.DB {
	db, err := dbcon.OpenSQLite(p.env.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return db
}
