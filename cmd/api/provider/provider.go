package provider

import (
	"database/sql"
	"log"

	"github.com/caarlos0/env/v11"

	"trama/internal/api/config"
	"trama/pkg/dbcon"
)

type Provider struct {
	cfg *config.Config
	db  *sql.DB
}

type envConfig struct {
	Port         string `env:"PORT,required"`
	GinMode      string `env:"GIN_MODE,required"`
	DatabasePath string `env:"DATABASE_PATH,required"`
}

func NewProvider() *Provider {
	var e envConfig
	if err := env.Parse(&e); err != nil {
		panic(err)
	}

	db, err := dbcon.OpenSQLite(e.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return &Provider{
		cfg: &config.Config{
			Port:         e.Port,
			GinMode:      e.GinMode,
			DatabasePath: e.DatabasePath,
		},
		db: db,
	}
}

func (p *Provider) Config() *config.Config {
	return p.cfg
}

func (p *Provider) DB() *sql.DB {
	return p.db
}
