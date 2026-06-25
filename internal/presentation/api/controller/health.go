package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/pkg/apiutil"
)

var (
	HealthPingSpec = apiutil.EndpointSpec{Verb: apiutil.GET, Path: "/ping", RequireAuth: false}
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Register(public gin.IRoutes, protected gin.IRoutes) {
	HealthPingSpec.RegisterOn(public, c.Ping)
}

// -----------------------------------------------------------------------------------
// Ping returns a simple ok response to indicate the service is running.
// @Summary      Health check
// @Description  Returns ok if the service is reachable
// @Tags         health
// @Success      200  {object}  map[string]string  "ok"
// @Router       /api/ping [get]
// -----------------------------------------------------------------------------------
func (c *HealthController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
