package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

const relayServiceSecretHeader = "X-Xphai-Relay-Secret"

// RelayServiceAuth restricts relay requests to trusted services when enabled
// in New API system settings.
func RelayServiceAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		enabled, expected := common.GetRelayServiceAuthConfig()
		if !enabled {
			c.Next()
			return
		}

		actual := c.GetHeader(relayServiceSecretHeader)
		if expected == "" || len(expected) != len(actual) ||
			subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"type":    "permission_error",
					"code":    "relay_service_access_denied",
					"message": "relay access is restricted to trusted services",
				},
			})
			return
		}

		c.Next()
	}
}
