package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trama/internal/api/apierror"
	"trama/internal/battlelog"
)

type matchService interface {
	List(ctx context.Context) ([]battlelog.MatchOutput, error)
	Create(ctx context.Context, in battlelog.CreateMatchInput) (battlelog.MatchOutput, error)
	Get(ctx context.Context, id uuid.UUID) (battlelog.MatchOutput, error)
	Update(ctx context.Context, in battlelog.UpdateMatchInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MatchController struct {
	service matchService
}

func NewMatchController(svc matchService) *MatchController {
	return &MatchController{service: svc}
}

func (c *MatchController) RegisterRoutes(rg *gin.RouterGroup) {
	matches := rg.Group("/matches")
	matches.GET("", c.list)
	matches.POST("", c.create)
	matches.GET("/:id", c.get)
	matches.PUT("/:id", c.update)
	matches.DELETE("/:id", c.delete)
}

// @Summary      List matches
// @Description  List all WH40k 11th matches
// @Tags         matches
// @Produce      json
// @Success      200  {array}   battlelog.MatchOutput
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/matches [get]
func (c *MatchController) list(ctx *gin.Context) {
	items, err := c.service.List(ctx.Request.Context())
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// @Summary      Create match
// @Description  Create a new WH40k 11th match
// @Tags         matches
// @Accept       json
// @Produce      json
// @Param        body  body      battlelog.CreateMatchInput  true  "Match data"
// @Success      201   {object}  battlelog.MatchOutput
// @Failure      400   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/matches [post]
func (c *MatchController) create(ctx *gin.Context) {
	var in battlelog.CreateMatchInput
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

// @Summary      Get match
// @Description  Get a match by ID
// @Tags         matches
// @Produce      json
// @Param        id   path      string  true  "Match ID"
// @Success      200  {object}  battlelog.MatchOutput
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/matches/{id} [get]
func (c *MatchController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid match id"))
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uid)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Update match
// @Description  Update a match by ID
// @Tags         matches
// @Accept       json
// @Param        id    path      string                    true  "Match ID"
// @Param        body  body      battlelog.UpdateMatchInput  true  "Match data"
// @Success      204   {object}  nil
// @Failure      400   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/matches/{id} [put]
func (c *MatchController) update(ctx *gin.Context) {
	var in battlelog.UpdateMatchInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}

	uid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid match id"))
		return
	}
	in.ID = uid

	if err := c.service.Update(ctx.Request.Context(), in); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Delete match
// @Description  Delete a match by ID
// @Tags         matches
// @Param        id   path      string  true  "Match ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/matches/{id} [delete]
func (c *MatchController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid match id"))
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uid); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
