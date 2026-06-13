package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/internal/core"
)

type GameSystemRepository interface {
	Create(gs core.GameSystem) (core.GameSystem, error)
	Get(id core.GameSystemID) (core.GameSystem, error)
	GetAll() ([]core.GameSystem, error)
	Update(gs core.GameSystem) error
	Delete(id core.GameSystemID) error
}

type EditionRepository interface {
	Create(e core.Edition) (core.Edition, error)
	Get(id core.EditionID) (core.Edition, error)
	GetAllByGameSystem(gsID core.GameSystemID) ([]core.Edition, error)
	Update(e core.Edition) error
	Delete(id core.EditionID) error
}

type FactionRepository interface {
	Create(f core.Faction) (core.Faction, error)
	Get(id core.FactionID) (core.Faction, error)
	GetAllByEdition(edID core.EditionID) ([]core.Faction, error)
	Update(f core.Faction) error
	Delete(id core.FactionID) error
}

type Handler struct {
	DB          *sql.DB
	GameSystems GameSystemRepository
	Editions    EditionRepository
	Factions    FactionRepository
}

func New(db *sql.DB, gs GameSystemRepository, ed EditionRepository, fac FactionRepository) *Handler {
	return &Handler{DB: db, GameSystems: gs, Editions: ed, Factions: fac}
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
