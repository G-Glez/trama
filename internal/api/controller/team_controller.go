package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trama/internal/api/apierror"
	"trama/internal/battlelog"
)

type teamService interface {
	List(ctx context.Context) ([]battlelog.TeamOutput, error)
	Create(ctx context.Context, in battlelog.CreateTeamInput) (battlelog.TeamOutput, error)
	Get(ctx context.Context, id uuid.UUID) (battlelog.TeamOutput, error)
	Update(ctx context.Context, in battlelog.UpdateTeamInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TeamController struct {
	service teamService
}

func NewTeamController(svc teamService) *TeamController {
	return &TeamController{service: svc}
}

func (c *TeamController) RegisterRoutes(rg *gin.RouterGroup) {
	teams := rg.Group("/teams")
	teams.GET("", c.list)
	teams.POST("", c.create)
	teams.GET("/:id", c.get)
	teams.PUT("/:id", c.update)
	teams.DELETE("/:id", c.delete)
}

// @Summary      List teams
// @Description  List all teams
// @Tags         teams
// @Produce      json
// @Success      200  {array}   battlelog.TeamOutput
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/teams [get]
func (c *TeamController) list(ctx *gin.Context) {
	items, err := c.service.List(ctx.Request.Context())
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// @Summary      Create team
// @Description  Create a new team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        body  body      battlelog.CreateTeamInput  true  "Team data"
// @Success      201   {object}  battlelog.TeamOutput
// @Failure      400   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/teams [post]
func (c *TeamController) create(ctx *gin.Context) {
	var in battlelog.CreateTeamInput
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

// @Summary      Get team
// @Description  Get a team by ID
// @Tags         teams
// @Produce      json
// @Param        id   path      string  true  "Team ID"
// @Success      200  {object}  battlelog.TeamOutput
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/teams/{id} [get]
func (c *TeamController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid team id"))
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uid)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Update team
// @Description  Update a team by ID
// @Tags         teams
// @Accept       json
// @Param        id    path      string                    true  "Team ID"
// @Param        body  body      battlelog.UpdateTeamInput  true  "Team data"
// @Success      204   {object}  nil
// @Failure      400   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/teams/{id} [put]
func (c *TeamController) update(ctx *gin.Context) {
	var in battlelog.UpdateTeamInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}

	uid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid team id"))
		return
	}
	in.ID = uid

	if err := c.service.Update(ctx.Request.Context(), in); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Delete team
// @Description  Delete a team by ID
// @Tags         teams
// @Param        id   path      string  true  "Team ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/teams/{id} [delete]
func (c *TeamController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid team id"))
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uid); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
