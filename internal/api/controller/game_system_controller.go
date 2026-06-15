package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trama/internal/api/apierror"
	"trama/internal/core"
)

type gameSystemService interface {
	GetAll(ctx context.Context) ([]core.GameSystemOutput, error)
	Get(ctx context.Context, id uuid.UUID) (core.GameSystemOutput, error)
	Create(ctx context.Context, in core.CreateGameSystemInput) (core.GameSystemOutput, error)
	Update(ctx context.Context, in core.UpdateGameSystemInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GameSystemController struct {
	service gameSystemService
}

func NewGameSystemController(svc gameSystemService) *GameSystemController {
	return &GameSystemController{service: svc}
}

func (c *GameSystemController) RegisterRoutes(rg *gin.RouterGroup) {
	gs := rg.Group("/game-systems")
	gs.GET("", c.list)
	gs.POST("", c.create)
	gs.GET("/:id", c.get)
	gs.PUT("/:id", c.update)
	gs.DELETE("/:id", c.delete)
}

// @Summary      List game systems
// @Description  Get all game systems
// @Tags         game-systems
// @Produce      json
// @Success      200  {array}   core.GameSystemOutput
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/game-systems [get]
func (c *GameSystemController) list(ctx *gin.Context) {
	items, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// @Summary      Create game system
// @Description  Create a new game system
// @Tags         game-systems
// @Accept       json
// @Produce      json
// @Param        body  body      core.CreateGameSystemInput  true  "Game system data"
// @Success      201   {object}  core.GameSystemOutput
// @Failure      400   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/game-systems [post]
func (c *GameSystemController) create(ctx *gin.Context) {
	var in core.CreateGameSystemInput
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

// @Summary      Get game system
// @Description  Get a game system by ID
// @Tags         game-systems
// @Produce      json
// @Param        id   path      string  true  "Game system ID"
// @Success      200  {object}  core.GameSystemOutput
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/game-systems/{id} [get]
func (c *GameSystemController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid game system id"))
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uid)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Update game system
// @Description  Update a game system by ID
// @Tags         game-systems
// @Accept       json
// @Param        id    path      string                    true  "Game system ID"
// @Param        body  body      core.UpdateGameSystemInput  true  "Game system data"
// @Success      204   {object}  nil
// @Failure      400   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/game-systems/{id} [put]
func (c *GameSystemController) update(ctx *gin.Context) {
	var in core.UpdateGameSystemInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}

	uid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid game system id"))
		return
	}
	in.ID = uid

	if err := c.service.Update(ctx.Request.Context(), in); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Delete game system
// @Description  Delete a game system by ID
// @Tags         game-systems
// @Param        id   path      string  true  "Game system ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/game-systems/{id} [delete]
func (c *GameSystemController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid game system id"))
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uid); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
