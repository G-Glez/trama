package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trama/internal/api/apierror"
	"trama/internal/core"
)

type factionService interface {
	GetAllByEdition(ctx context.Context, edID uuid.UUID) ([]core.FactionOutput, error)
	Create(ctx context.Context, in core.CreateFactionInput) (core.FactionOutput, error)
	Get(ctx context.Context, id uuid.UUID) (core.FactionOutput, error)
	Update(ctx context.Context, in core.UpdateFactionInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FactionController struct {
	service factionService
}

func NewFactionController(svc factionService) *FactionController {
	return &FactionController{service: svc}
}

func (c *FactionController) RegisterRoutes(rg *gin.RouterGroup) {
	factions := rg.Group("/factions")
	factions.GET("", c.list)
	factions.POST("", c.create)
	factions.GET("/:id", c.get)
	factions.PUT("/:id", c.update)
	factions.DELETE("/:id", c.delete)
}

// @Summary      List factions
// @Description  Get all factions for an edition
// @Tags         factions
// @Produce      json
// @Param        edition_id  query     string  true  "Edition ID"
// @Success      200         {array}   core.FactionOutput
// @Failure      400         {object}  apierror.Error
// @Failure      500         {object}  apierror.Error
// @Router       /api/v1/factions [get]
func (c *FactionController) list(ctx *gin.Context) {
	var req struct {
		EditionID uuid.UUID `form:"edition_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}
	items, err := c.service.GetAllByEdition(ctx.Request.Context(), req.EditionID)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// @Summary      Create faction
// @Description  Create a new faction for an edition
// @Tags         factions
// @Accept       json
// @Produce      json
// @Param        body  body      core.CreateFactionInput  true  "Faction data"
// @Success      201   {object}  core.FactionOutput
// @Failure      400   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/factions [post]
func (c *FactionController) create(ctx *gin.Context) {
	var in core.CreateFactionInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}
	out, err := c.service.Create(ctx.Request.Context(), in)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// @Summary      Get faction
// @Description  Get a faction by ID
// @Tags         factions
// @Produce      json
// @Param        id   path      string  true  "Faction ID"
// @Success      200  {object}  core.FactionOutput
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/factions/{id} [get]
func (c *FactionController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid faction id"))
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uid)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Update faction
// @Description  Update a faction by ID
// @Tags         factions
// @Accept       json
// @Param        id    path      string                    true  "Faction ID"
// @Param        body  body      core.UpdateFactionInput   true  "Faction data"
// @Success      204   {object}  nil
// @Failure      400   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/factions/{id} [put]
func (c *FactionController) update(ctx *gin.Context) {
	var in core.UpdateFactionInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}

	uid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid faction id"))
		return
	}
	in.ID = uid

	if err := c.service.Update(ctx.Request.Context(), in); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Delete faction
// @Description  Delete a faction by ID
// @Tags         factions
// @Param        id   path      string  true  "Faction ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/factions/{id} [delete]
func (c *FactionController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid faction id"))
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uid); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
