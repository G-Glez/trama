package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trama/internal/api/apierror"
	"trama/internal/core"
)

type editionService interface {
	GetAllByGameSystem(ctx context.Context, gsID uuid.UUID) ([]core.EditionOutput, error)
	Create(ctx context.Context, in core.CreateEditionInput) (core.EditionOutput, error)
	Get(ctx context.Context, id uuid.UUID) (core.EditionOutput, error)
	Update(ctx context.Context, in core.UpdateEditionInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EditionController struct {
	service editionService
}

func NewEditionController(svc editionService) *EditionController {
	return &EditionController{service: svc}
}

func (c *EditionController) RegisterRoutes(rg *gin.RouterGroup) {
	ed := rg.Group("/editions")
	ed.GET("", c.list)
	ed.POST("", c.create)
	ed.GET("/:id", c.get)
	ed.PUT("/:id", c.update)
	ed.DELETE("/:id", c.delete)
}

// @Summary      List editions
// @Description  Get all editions for a game system
// @Tags         editions
// @Produce      json
// @Param        game_system_id  query     string  true  "Game system ID"
// @Success      200             {array}   core.EditionOutput
// @Failure      400             {object}  apierror.Error
// @Failure      500             {object}  apierror.Error
// @Router       /api/v1/editions [get]
func (c *EditionController) list(ctx *gin.Context) {
	var req struct {
		GameSystemID uuid.UUID `form:"game_system_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}
	items, err := c.service.GetAllByGameSystem(ctx.Request.Context(), req.GameSystemID)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// @Summary      Create edition
// @Description  Create a new edition for a game system
// @Tags         editions
// @Accept       json
// @Produce      json
// @Param        body  body      core.CreateEditionInput  true  "Edition data"
// @Success      201   {object}  core.EditionOutput
// @Failure      400   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/editions [post]
func (c *EditionController) create(ctx *gin.Context) {
	var in core.CreateEditionInput
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

// @Summary      Get edition
// @Description  Get an edition by ID
// @Tags         editions
// @Produce      json
// @Param        id   path      string  true  "Edition ID"
// @Success      200  {object}  core.EditionOutput
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/editions/{id} [get]
func (c *EditionController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid edition id"))
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uid)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Update edition
// @Description  Update an edition by ID
// @Tags         editions
// @Accept       json
// @Param        id    path      string                    true  "Edition ID"
// @Param        body  body      core.UpdateEditionInput   true  "Edition data"
// @Success      204   {object}  nil
// @Failure      400   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/editions/{id} [put]
func (c *EditionController) update(ctx *gin.Context) {
	var in core.UpdateEditionInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}

	uid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid edition id"))
		return
	}
	in.ID = uid

	if err := c.service.Update(ctx.Request.Context(), in); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Delete edition
// @Description  Delete an edition by ID
// @Tags         editions
// @Param        id   path      string  true  "Edition ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/editions/{id} [delete]
func (c *EditionController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid edition id"))
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uid); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
