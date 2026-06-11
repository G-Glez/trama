package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"trama/config"
	_ "trama/docs"
)

// @title           TRAMA API
// @version         1.0
// @description     Tournament Records and Metrics Assistant
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	r.GET("/ping", pingHandler)
	r.GET("/hola", holaHandler)

	if cfg.GinMode == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}

// @Summary      Health check
// @Description  Returns pong
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// @Summary      Say hello
// @Description  Returns a greeting
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /hola [get]
func holaHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hola mundo",
	})
}
