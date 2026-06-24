package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"trama/internal/domain/core/principal"
	"trama/internal/presentation/api/auth"
)

type AuthService interface {
	ValidateToken(token string) (auth.Claims, error)
}

// -----------------------------------------------------------------------------------
// AuthMiddleware validates a Bearer token from the Authorization header and
// injects a principal.Principal into the request context.
// Aborts with 401 if the token is missing, malformed, or invalid.
// -----------------------------------------------------------------------------------
func AuthMiddleware(svc AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := svc.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		ppl := principal.New(claims.UserID, []string{}, []string{})
		ctx := principal.WithPrincipal(c.Request.Context(), ppl)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
