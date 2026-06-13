package provider

import (
	"database/sql"
	"log"

	"github.com/caarlos0/env/v11"

	"trama/internal/cli"
	"trama/pkg/dbcon"
)

type envConfig struct {
	DatabasePath string `env:"DATABASE_PATH"`
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

func (p *Provider) Cli(db *sql.DB) *cli.CLI {
	return cli.New(db)
}

func (p *Provider) DBCon() *sql.DB {
	db, err := dbcon.OpenSQLite(p.env.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	return db
}
