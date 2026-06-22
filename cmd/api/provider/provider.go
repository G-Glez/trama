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
	"trama/internal/battlelog"
	battleloggen "trama/internal/gen/battlelog"
	coregen "trama/internal/gen/core"
	usergen "trama/internal/gen/user"
	"trama/internal/core"
	"trama/internal/user"
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
	qCore   *coregen.Queries
	qUser   *usergen.Queries
	qBattle *battleloggen.Queries
	facRepo *core.FactionSQLRepository
	userRepo *user.UserSQLRepository
	matchRepo *battlelog.MatchSQLRepository
	teamRepo *battlelog.TeamSQLRepository
	tournRepo *battlelog.TournamentSQLRepository
	facSvc  *core.FactionService
	userSvc *user.UserService
	matchSvc *battlelog.MatchService
	teamSvc *battlelog.TeamService
	tournSvc *battlelog.TournamentService
	facCtrl *controller.FactionController
	userCtrl *controller.UserController
	matchCtrl *controller.MatchController
	teamCtrl *controller.TeamController
	tournCtrl *controller.TournamentController
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
	if p.qCore != nil {
		return
	}

	p.qCore = coregen.New(p.db)
	p.qUser = usergen.New(p.db)
	p.qBattle = battleloggen.New(p.db)
}

func (p *Provider) provisionRepos() {
	if p.facRepo != nil {
		return
	}

	p.facRepo = core.NewFactionRepository(p.qCore)
	p.userRepo = user.NewUserRepository(p.qUser)
	p.matchRepo = battlelog.NewMatchRepository(p.qBattle)
	p.teamRepo = battlelog.NewTeamRepository(p.qBattle)
	p.tournRepo = battlelog.NewTournamentRepository(p.qBattle)
}

func (p *Provider) provisionServices() {
	if p.facSvc != nil {
		return
	}

	p.facSvc = core.NewFactionService(p.facRepo)
	p.userSvc = user.NewUserService(p.userRepo)
	p.matchSvc = battlelog.NewMatchService(p.matchRepo)
	p.teamSvc = battlelog.NewTeamService(p.teamRepo)
	p.tournSvc = battlelog.NewTournamentService(p.tournRepo)
}

func (p *Provider) provisionGin() {
	if p.gin != nil {
		return
	}

	gin.SetMode(p.env.GinMode)
	p.gin = gin.Default()
}

func (p *Provider) provisionControllers() {
	if p.facCtrl != nil {
		return
	}

	p.facCtrl = controller.NewFactionController(p.facSvc)
	p.userCtrl = controller.NewUserController(p.userSvc)
	p.matchCtrl = controller.NewMatchController(p.matchSvc)
	p.teamCtrl = controller.NewTeamController(p.teamSvc)
	p.tournCtrl = controller.NewTournamentController(p.tournSvc)
}

func (p *Provider) provisionRouter() {
	if p.router != nil {
		return
	}

	r := api.NewRouter(p.gin, p.Config())
	r.WithController(p.facCtrl)
	r.WithController(p.userCtrl)
	r.WithController(p.matchCtrl)
	r.WithController(p.teamCtrl)
	r.WithController(p.tournCtrl)
	r.Setup()

	p.router = r
}
