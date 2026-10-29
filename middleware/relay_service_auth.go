package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

const relayServiceSecretHeader = "X-Xphai-Relay-Secret"

// RelayServiceAuth restricts relay requests to trusted services when enabled
// in New API system settings.
func RelayServiceAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if requireRelayServiceSecret(c) {
			c.Next()
		}
	}
}

func requireRelayServiceSecret(c *gin.Context) bool {
	enabled, expected := common.GetRelayServiceAuthConfig()
	if !enabled {
		return true
	}
	var values []string
	for name, entries := range c.Request.Header {
		if strings.EqualFold(name, relayServiceSecretHeader) {
			values = append(values, entries...)
		}
	}
	if expected != "" && len(values) == 1 && len(expected) == len(values[0]) &&
		subtle.ConstantTimeCompare([]byte(expected), []byte(values[0])) == 1 {
		return true
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error": gin.H{
			"type":    "permission_error",
			"code":    "relay_service_access_denied",
			"message": "relay access is restricted to trusted services",
		},
	})
	return false
}
