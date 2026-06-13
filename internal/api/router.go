// @title           TRAMA API
// @version         1.0
// @description     Tournament Records and Metrics Assistant
// @host            localhost:8080
// @BasePath        /

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"trama/internal/api/config"
)

type Controller interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

type Router struct {
	engine      *gin.Engine
	controllers []Controller
}

func NewRouter(engine *gin.Engine, cfg *config.Config) *Router {
	engine.GET("/ping", Ping)

	if cfg.GinMode == "debug" {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return &Router{engine: engine}
}

func (rt *Router) WithController(c Controller) *Router {
	rt.controllers = append(rt.controllers, c)
	return rt
}

func (rt *Router) Setup() {
	api := rt.engine.Group("/api/v1")
	for _, c := range rt.controllers {
		c.RegisterRoutes(api)
	}
}

func (rt *Router) Run(addr string) error {
	return rt.engine.Run(addr)
}

// @Summary      Health check
// @Description  Returns pong
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
