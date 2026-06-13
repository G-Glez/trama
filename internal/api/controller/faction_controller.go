package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/internal/core"
)

type FactionController struct {
	service *core.FactionService
}

func NewFactionController(svc *core.FactionService) *FactionController {
	return &FactionController{service: svc}
}

func (c *FactionController) RegisterRoutes(rg *gin.RouterGroup) {
	factions := rg.Group("/factions")
	factions.GET("", c.List)
	factions.POST("", c.Create)
	factions.GET("/:id", c.Get)
	factions.PUT("/:id", c.Update)
	factions.DELETE("/:id", c.Delete)
}

func (c *FactionController) List(ctx *gin.Context) {
	edID := ctx.Query("edition_id")
	if edID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "edition_id query parameter is required"})
		return
	}
	items, err := c.service.GetAllByEdition(edID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *FactionController) Create(ctx *gin.Context) {
	var in core.CreateFactionInput
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

func (c *FactionController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	out, err := c.service.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, out)
}

func (c *FactionController) Update(ctx *gin.Context) {
	var in core.UpdateFactionInput
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

func (c *FactionController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
