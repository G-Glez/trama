package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/internal/presentation/api/auth"
	"trama/pkg/apiutil"
)

var (
	AuthRegisterSpec = apiutil.EndpointSpec{Verb: apiutil.POST, Path: "/auth/register", RequireAuth: false}
	AuthLoginSpec    = apiutil.EndpointSpec{Verb: apiutil.POST, Path: "/auth/login", RequireAuth: false}
)

type AuthController struct {
	svc *auth.Service
}

func NewAuthController(svc *auth.Service) *AuthController {
	return &AuthController{svc: svc}
}

func (c *AuthController) Register(public gin.IRoutes, protected gin.IRoutes) {
	AuthRegisterSpec.RegisterOn(public, c.RegisterHandler)
	AuthLoginSpec.RegisterOn(public, c.LoginHandler)
}

// -----------------------------------------------------------------------------------
// RegisterHandler creates a new user account.
// @Summary      Register a new user
// @Description  Creates a new user account with the given username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration credentials"
// @Success      201  {object}
// @Failure      400  {object}  errorhandler.ErrorResponse  "validation error"
// @Failure      409  {object}  errorhandler.ErrorResponse  "username already exists"
// @Failure      500  {object}  errorhandler.ErrorResponse  "internal error"
// @Router       /api/auth/register [post]
// -----------------------------------------------------------------------------------
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"john"`
	Password string `json:"password" binding:"required,min=6,max=100" example:"secret123"`
}

func (c *AuthController) RegisterHandler(ctx *gin.Context) {
	var req RegisterRequest
	if ok := body(ctx, &req); !ok {
		return
	}

	if err := c.svc.Register(ctx.Request.Context(), auth.RegisterInput(req)); err != nil {
		fail(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

// -----------------------------------------------------------------------------------
// LoginHandler authenticates a user and returns a JWT token.
// @Summary      Login
// @Description  Authenticates with username and password, returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  errorhandler.ErrorResponse  "validation error"
// @Failure      401  {object}  errorhandler.ErrorResponse  "invalid credentials"
// @Failure      500  {object}  errorhandler.ErrorResponse  "internal error"
// @Router       /api/auth/login [post]
// -----------------------------------------------------------------------------------
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john"`
	Password string `json:"password" binding:"required" example:"secret123"`
}
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

func (c *AuthController) LoginHandler(ctx *gin.Context) {
	var req LoginRequest
	if ok := body(ctx, &req); !ok {
		return
	}

	resp, err := c.svc.Login(ctx.Request.Context(), auth.LoginInput(req))
	if err != nil {
		fail(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse(resp))
}
