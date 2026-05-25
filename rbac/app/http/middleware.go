package http

import (
	"crypto/subtle"
	nethttp "net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SecretAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{
				"error": "invalid authorization scheme",
			})
			return
		}

		providedSecret := strings.TrimPrefix(authHeader, prefix)

		if subtle.ConstantTimeCompare([]byte(providedSecret), []byte(secret)) != 1 {
			c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{
				"error": "invalid secret",
			})
			return
		}

		c.Next()
	}
}
