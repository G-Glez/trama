package provider

import (
	"database/sql"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"

	_ "trama/docs"
	"trama/internal/api"
	"trama/internal/api/config"
	"trama/internal/api/controller"
	"trama/internal/core"
	coregen "trama/internal/gen/core"
	"trama/pkg/dbcon"
)

type envConfig struct {
	Port         string `env:"PORT,required"`
	GinMode      string `env:"GIN_MODE,required"`
	DatabasePath string `env:"DATABASE_PATH,required"`
}

type Provider struct {
	env     envConfig
	db      *sql.DB
	queries *coregen.Queries
	gsRepo  *core.GameSystemSQLRepository
	edRepo  *core.EditionSQLRepository
	facRepo *core.FactionSQLRepository
	gsSvc   *core.GameSystemService
	edSvc   *core.EditionService
	facSvc  *core.FactionService
	gsCtrl  *controller.GameSystemController
	edCtrl  *controller.EditionController
	facCtrl *controller.FactionController
	gin     *gin.Engine
	router  *api.Router
}

func NewProvisionedProvider() *Provider {
	var e envConfig
	if err := env.Parse(&e); err != nil {
		panic(err)
	}

	p := &Provider{env: e}

	p.provisionDB()
	p.provisionQueries()
	p.provisionRepos()
	p.provisionServices()
	p.provisionControllers()
	p.provisionGin()
	p.provisionRouter()

	return p
}

func (p *Provider) Config() *config.Config {
	return &config.Config{
		Port:         p.env.Port,
		GinMode:      p.env.GinMode,
		DatabasePath: p.env.DatabasePath,
	}
}

func (p *Provider) DB() *sql.DB {
	return p.db
}

func (p *Provider) Router() *api.Router {
	return p.router
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

func (p *Provider) provisionQueries() {
	if p.queries != nil {
		return
	}

	p.queries = coregen.New(p.db)
}

func (p *Provider) provisionRepos() {
	if p.gsRepo != nil {
		return
	}

	p.gsRepo = core.NewGameSystemRepository(p.queries)
	p.edRepo = core.NewEditionRepository(p.queries)
	p.facRepo = core.NewFactionRepository(p.queries)
}

func (p *Provider) provisionServices() {
	if p.gsSvc != nil {
		return
	}

	p.gsSvc = core.NewGameSystemService(p.gsRepo)
	p.edSvc = core.NewEditionService(p.edRepo)
	p.facSvc = core.NewFactionService(p.facRepo)
}

func (p *Provider) provisionGin() {
	if p.gin != nil {
		return
	}

	gin.SetMode(p.env.GinMode)
	p.gin = gin.Default()
}

func (p *Provider) provisionControllers() {
	if p.gsCtrl != nil {
		return
	}

	p.gsCtrl = controller.NewGameSystemController(p.gsSvc)
	p.edCtrl = controller.NewEditionController(p.edSvc)
	p.facCtrl = controller.NewFactionController(p.facSvc)
}

func (p *Provider) provisionRouter() {
	if p.router != nil {
		return
	}

	r := api.NewRouter(p.gin, p.Config())
	r.WithController(p.gsCtrl)
	r.WithController(p.edCtrl)
	r.WithController(p.facCtrl)
	r.Setup()

	p.router = r
}
