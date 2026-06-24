// @title           TRAMA API
// @version         1.0
// @description     Tournament Records and Metrics Assistant
// @BasePath        /

package api

import "github.com/gin-gonic/gin"

type Controller interface {
	Register(public gin.IRoutes, protected gin.IRoutes)
}

type Router struct {
	controllers []Controller
}

func NewRouter(controllers ...Controller) *Router {
	return &Router{controllers: controllers}
}

func (rt *Router) Setup(
	r *gin.Engine,
	authMiddleware gin.HandlerFunc,
) {
	api := r.Group("/api")

	public := api.Group("")
	protected := api.Group("")
	protected.Use(authMiddleware)

	for _, c := range rt.controllers {
		c.Register(public, protected)
	}
}
