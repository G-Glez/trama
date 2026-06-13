package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/internal/core"
)

type EditionController struct {
	service *core.EditionService
}

func NewEditionController(svc *core.EditionService) *EditionController {
	return &EditionController{service: svc}
}

func (c *EditionController) RegisterRoutes(rg *gin.RouterGroup) {
	ed := rg.Group("/editions")
	ed.GET("", c.List)
	ed.POST("", c.Create)
	ed.GET("/:id", c.Get)
	ed.PUT("/:id", c.Update)
	ed.DELETE("/:id", c.Delete)
}

func (c *EditionController) List(ctx *gin.Context) {
	gsID := ctx.Query("game_system_id")
	if gsID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "game_system_id query parameter is required"})
		return
	}
	items, err := c.service.GetAllByGameSystem(gsID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *EditionController) Create(ctx *gin.Context) {
	var in core.CreateEditionInput
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

func (c *EditionController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	out, err := c.service.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, out)
}

func (c *EditionController) Update(ctx *gin.Context) {
	var in core.UpdateEditionInput
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

func (c *EditionController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
