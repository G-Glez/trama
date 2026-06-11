package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"trama/config"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
