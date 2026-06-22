package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trama/internal/api/apierror"
	"trama/internal/user"
)

type userService interface {
	Register(ctx context.Context, in user.RegisterInput) (user.UserOutput, error)
	Login(ctx context.Context, in user.LoginInput) (user.UserOutput, error)
	Get(ctx context.Context, id uuid.UUID) (user.UserOutput, error)
	Update(ctx context.Context, in user.UpdateUserInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserController struct {
	service userService
}

func NewUserController(svc userService) *UserController {
	return &UserController{service: svc}
}

func (c *UserController) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.POST("/register", c.register)
	users.POST("/login", c.login)
	users.GET("/:id", c.get)
	users.PUT("/:id", c.update)
	users.DELETE("/:id", c.delete)
}

// @Summary      Register user
// @Description  Register a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      user.RegisterInput  true  "User data"
// @Success      201   {object}  user.UserOutput
// @Failure      400   {object}  apierror.Error
// @Failure      409   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/users/register [post]
func (c *UserController) register(ctx *gin.Context) {
	var in user.RegisterInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}
	out, err := c.service.Register(ctx.Request.Context(), in)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// @Summary      Login
// @Description  Authenticate a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      user.LoginInput  true  "Credentials"
// @Success      200   {object}  user.UserOutput
// @Failure      400   {object}  apierror.Error
// @Failure      401   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/users/login [post]
func (c *UserController) login(ctx *gin.Context) {
	var in user.LoginInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}
	out, err := c.service.Login(ctx.Request.Context(), in)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Get user
// @Description  Get a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  user.UserOutput
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/users/{id} [get]
func (c *UserController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid user id"))
		return
	}
	out, err := c.service.Get(ctx.Request.Context(), uid)
	if err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary      Update user
// @Description  Update a user by ID
// @Tags         users
// @Accept       json
// @Param        id    path      string                  true  "User ID"
// @Param        body  body      user.UpdateUserInput    true  "User data"
// @Success      204   {object}  nil
// @Failure      400   {object}  apierror.Error
// @Failure      404   {object}  apierror.Error
// @Failure      500   {object}  apierror.Error
// @Router       /api/v1/users/{id} [put]
func (c *UserController) update(ctx *gin.Context) {
	var in user.UpdateUserInput
	if err := ctx.ShouldBindJSON(&in); err != nil {
		apierror.HandleError(ctx, apierror.BadRequest(err.Error()))
		return
	}

	uid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid user id"))
		return
	}
	in.ID = uid

	if err := c.service.Update(ctx.Request.Context(), in); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Param        id   path      string  true  "User ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /api/v1/users/{id} [delete]
func (c *UserController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		apierror.HandleError(ctx, apierror.BadRequest("invalid user id"))
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), uid); err != nil {
		apierror.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
