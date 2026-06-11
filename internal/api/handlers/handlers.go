package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// @Summary      Health check
// @Description  Returns pong
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// @Summary      Say hello
// @Description  Returns a greeting and stores an ID in SQLite
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /hola [get]
func (h *Handler) Hola(c *gin.Context) {
	result, err := h.DB.Exec("INSERT INTO greetings DEFAULT VALUES")
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
