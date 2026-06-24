package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"trama/internal/domain/core/principal"
)

// -----------------------------------------------------------------------------------
// AuthMiddleware returns a Gin handler that validates Bearer tokens and injects a
// Principal into the request context. The middleware is intended to be applied to
// protected route groups only — public endpoints should not pass through it.
// -----------------------------------------------------------------------------------
func AuthMiddleware(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := svc.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		p := principal.New(claims.Username, nil, nil)
		ctx := principal.WithPrincipal(c.Request.Context(), p)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
