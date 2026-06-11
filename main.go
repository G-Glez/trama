package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"trama/config"
	_ "trama/docs"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// @title           TRAMA API
// @version         1.0
// @description     Tournament Records and Metrics Assistant
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.Load()

	var err error
	db, err = sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS greetings (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

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
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// @Summary      Say hello
// @Description  Returns a greeting and stores an ID in SQLite
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /hola [get]
func holaHandler(c *gin.Context) {
	result, err := db.Exec("INSERT INTO greetings DEFAULT VALUES")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hola mundo",
		"id":      id,
	})
}
