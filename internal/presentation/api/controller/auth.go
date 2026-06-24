package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/internal/presentation/api/auth"
	"trama/pkg/apidef"
)

var (
	AuthRegisterSpec = apidef.EndpointSpec{Verb: apidef.POST, Path: "/auth/register", RequireAuth: false}
	AuthLoginSpec    = apidef.EndpointSpec{Verb: apidef.POST, Path: "/auth/login", RequireAuth: false}
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

func (c *AuthController) RegisterHandler(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6,max=100"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.Register(ctx.Request.Context(), auth.RegisterInput(req)); err != nil {
		if err == auth.ErrUserExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (c *AuthController) LoginHandler(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.svc.Login(ctx.Request.Context(), auth.LoginInput(req))
	if err != nil {
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidPassword {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
