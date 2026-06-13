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
	env   envConfig
	dbCon *sql.DB
	cli   *cli.CLI
}

func NewProvisionedProvider() *Provider {
	var e envConfig
	if err := env.Parse(&e); err != nil {
		panic(err)
	}

	p := &Provider{env: e}

	p.provisionDB()
	p.provisionCLI()

	return p
}

func (p *Provider) CLI() *cli.CLI {
	return p.cli
}

func (p *Provider) provisionDB() {
	if p.dbCon != nil {
		return
	}

	db, err := dbcon.OpenSQLite(p.env.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	p.dbCon = db
}

func (p *Provider) provisionCLI() {
	if p.cli == nil {
		p.provisionDB()
	}

	p.cli = cli.New(p.dbCon)
}
