package provider

import (
	"trama/internal/presentation/api"
	"trama/internal/presentation/api/auth"
	"trama/internal/presentation/api/controller"

	"github.com/gin-gonic/gin"
)

type ApiProvider struct {
	jwt              *auth.JWT
	gin              *gin.Engine
	authService      *auth.Service
	authController   *controller.AuthController
	healthController *controller.HealthController
	authMiddleware   gin.HandlerFunc
	router           *api.Router
}

func (p *Provider) provisionApi() {
	p.provisionJWT()
	p.provisionAuthService()
	p.provisionAuthController()
	p.provisionHealthController()
	p.provisionAuthMiddleware()
	p.provisionRouter()
	p.provisionGin()
}

func (p *Provider) provisionJWT() {
	p.jwt = auth.NewJWT([]byte(p.env.JWTSecret))
}

func (p *Provider) provisionAuthService() {
	p.authService = auth.NewService(p.userRepository, p.jwt)
}

func (p *Provider) provisionAuthController() {
	p.authController = controller.NewAuthController(p.authService)
}

func (p *Provider) provisionHealthController() {
	p.healthController = controller.NewHealthController()
}

func (p *Provider) provisionAuthMiddleware() {
	p.authMiddleware = auth.AuthMiddleware(p.authService)
}

func (p *Provider) provisionRouter() {
	p.router = api.NewRouter(
		p.authController,
		p.healthController,
	)
}

func (p *Provider) provisionGin() {
	ginMode := "release"
	if p.env.Env == "dev" {
		ginMode = "debug"
	}

	gin.SetMode(ginMode)
	p.gin = gin.Default()
}
