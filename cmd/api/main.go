package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"trama/cmd/api/provider"
	_ "trama/docs"
	"trama/internal/api/handlers"
	"trama/internal/core"
	coregen "trama/internal/gen/core"
)

// @title           TRAMA API
// @version         1.0
// @description     Tournament Records and Metrics Assistant
// @host            localhost:8080
// @BasePath        /
func main() {
	p := provider.NewProvider()
	defer p.DB().Close()

	cfg := p.Config()
	db := p.DB()

	q := coregen.New(db)
	h := handlers.New(db,
		core.NewGameSystemRepository(q),
		core.NewEditionRepository(q),
		core.NewFactionRepository(q),
	)

	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	r.GET("/ping", h.Ping)

	if cfg.GinMode == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
