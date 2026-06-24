package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/pkg/apidef"
)

var (
	HealthPingSpec      = apidef.EndpointSpec{Verb: apidef.GET, Path: "/ping", RequireAuth: false}
	HealthEndpointSpecs = []apidef.EndpointSpec{HealthPingSpec}
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Register(public gin.IRoutes, protected gin.IRoutes) {
	HealthPingSpec.RegisterOn(public, c.Ping)
}

func (c *HealthController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
