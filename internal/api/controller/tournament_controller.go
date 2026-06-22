package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trama/internal/api/apierror"
	"trama/internal/battlelog"
)

type tournamentService interface {
	List(ctx context.Context) ([]battlelog.TournamentOutput, error)
	Create(ctx context.Context, in battlelog.CreateTournamentInput) (battlelog.TournamentOutput, error)
	Get(ctx context.Context, id uuid.UUID) (battlelog.TournamentOutput, error)
	Update(ctx context.Context, in battlelog.UpdateTournamentInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TournamentController struct {
	service tournamentService
}

func NewTournamentController(svc tournamentService) *TournamentController {
	return &TournamentController{service: svc}
}

func (c *TournamentController) RegisterRoutes(rg *gin.RouterGroup) {
	tournaments := rg.Group("/tournaments")
	tournaments.GET("", c.list)
	tournaments.POST("", c.create)
	tournaments.GET("/:id", c.get)
	tournaments.PUT("/:id", c.update)
	tournaments.DELETE("/:id", c.delete)
}

// @Summary      List tournaments
// @Description  List all tournaments
// @Tags         tournaments
// @Produce      json
// @Success      200  {array}   battlelog.TournamentOutput
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/tournaments [get]
func (c *TournamentController) list(ctx *gin.Context) {
	items, err := c.service.List(ctx.Request.Context())
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// @Summary      Create tournament
// @Description  Create a new tournament
// @Tags         tournaments
// @Accept       json
// @Produce      json
// @Param        body  body      battlelog.CreateTournamentInput  true  "Tournament data"
// @Success      201   {object}  battlelog.TournamentOutput
// @Failure      400   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/tournaments [post]
func (c *TournamentController) create(ctx *gin.Context) {
	var in battlelog.CreateTournamentInput
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

// @Summary      Get tournament
// @Description  Get a tournament by ID
// @Tags         tournaments
// @Produce      json
// @Param        id   path      string  true  "Tournament ID"
// @Success      200  {object}  battlelog.TournamentOutput
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/tournaments/{id} [get]
func (c *TournamentController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid tournament id"))
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uid)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Update tournament
// @Description  Update a tournament by ID
// @Tags         tournaments
// @Accept       json
// @Param        id    path      string                         true  "Tournament ID"
// @Param        body  body      battlelog.UpdateTournamentInput  true  "Tournament data"
// @Success      204   {object}  nil
// @Failure      400   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/tournaments/{id} [put]
func (c *TournamentController) update(ctx *gin.Context) {
	var in battlelog.UpdateTournamentInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}

	uid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid tournament id"))
		return
	}
	in.ID = uid

	if err := c.service.Update(ctx.Request.Context(), in); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Delete tournament
// @Description  Delete a tournament by ID
// @Tags         tournaments
// @Param        id   path      string  true  "Tournament ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/tournaments/{id} [delete]
func (c *TournamentController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid tournament id"))
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uid); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
