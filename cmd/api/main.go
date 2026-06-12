package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"trama/internal/api/config"
	"trama/internal/api/database"
	"trama/internal/api/handlers"
	"trama/internal/core"
	coregen "trama/internal/gen/core"
	_ "trama/docs"
)

// @title           TRAMA API
// @version         1.0
// @description     Tournament Records and Metrics Assistant
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	q := coregen.New(db)
	h := handlers.New(db,
		core.NewGameSystemRepository(q),
		core.NewEditionRepository(q),
		core.NewFactionRepository(q),
	)

	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	r.GET("/ping", h.Ping)
	r.GET("/hola", h.Hola)

	if cfg.GinMode == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
