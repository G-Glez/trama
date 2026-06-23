package provider

import (
	"database/sql"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"

	"trama/internal/api"
	"trama/pkg/dbcon"
)

type envConfig struct {
	GinMode      string `env:"GIN_MODE,required"`
	DatabasePath string `env:"DATABASE_PATH,required"`
}

type Provider struct {
	env envConfig
	db  *sql.DB
	gin *gin.Engine
}

func NewProvisionedProvider() *Provider {
	var e envConfig
	if err := env.Parse(&e); err != nil {
		panic(err)
	}

	p := &Provider{env: e}

	p.provisionDB()
	p.provisionGin()

	return p
}

func (p *Provider) Gin() *gin.Engine {
	return p.gin
}

func (p *Provider) DB() *sql.DB {
	return p.db
}

func (p *Provider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *Provider) provisionDB() {
	if p.db != nil {
		return
	}

	db, err := dbcon.OpenSQLite(p.env.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	p.db = db
}

func (p *Provider) provisionGin() {
	if p.gin != nil {
		return
	}

	gin.SetMode(p.env.GinMode)
	p.gin = gin.Default()
	p.gin.GET("/api/ping", api.Ping)
}
