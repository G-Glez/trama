package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/internal/core"
)

type GameSystemController struct {
	service *core.GameSystemService
}

func NewGameSystemController(svc *core.GameSystemService) *GameSystemController {
	return &GameSystemController{service: svc}
}

func (c *GameSystemController) RegisterRoutes(rg *gin.RouterGroup) {
	gs := rg.Group("/game-systems")
	gs.GET("", c.List)
	gs.POST("", c.Create)
	gs.GET("/:id", c.Get)
	gs.PUT("/:id", c.Update)
	gs.DELETE("/:id", c.Delete)
}

func (c *GameSystemController) List(ctx *gin.Context) {
	items, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *GameSystemController) Create(ctx *gin.Context) {
	var in core.CreateGameSystemInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := c.service.Create(in)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

func (c *GameSystemController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	out, err := c.service.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, out)
}

func (c *GameSystemController) Update(ctx *gin.Context) {
	var in core.UpdateGameSystemInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	in.ID = ctx.Param("id")
	if err := c.service.Update(in); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *GameSystemController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
