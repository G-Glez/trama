// @title           TRAMA API
// @version         1.0
// @description     Tournament Records and Metrics Assistant
// @BasePath        /

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Health check
// @Description  Returns ok
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
