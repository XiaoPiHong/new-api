package controller

import (
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

type relayServiceAuthConfigResponse struct {
	Enabled          bool `json:"enabled"`
	SecretConfigured bool `json:"secret_configured"`
}

type relayServiceAuthUpdateRequest struct {
	Enabled *bool   `json:"enabled"`
	Secret  *string `json:"secret"`
}

func currentRelayServiceAuthResponse() relayServiceAuthConfigResponse {
	enabled, secret := common.GetRelayServiceAuthConfig()
	return relayServiceAuthConfigResponse{
		Enabled:          enabled,
		SecretConfigured: secret != "",
	}
}

// GetRelayServiceAuth returns only safe status fields. The configured secret
// is never serialized into an API response.
func GetRelayServiceAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    currentRelayServiceAuthResponse(),
	})
}

// UpdateRelayServiceAuth persists the relay protection state in New API's
// existing Option table and updates OptionMap immediately through the model
// option API. Omitting secret keeps the current secret for enable/disable
// operations; sending an empty string clears it.
func UpdateRelayServiceAuth(c *gin.Context) {
	var request relayServiceAuthUpdateRequest
	if err := common.DecodeJson(c.Request.Body, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if request.Enabled == nil && request.Secret == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "至少提供一个配置项",
		})
		return
	}

	currentEnabled, currentSecret := common.GetRelayServiceAuthConfig()
	enabled := currentEnabled
	secret := currentSecret
	if request.Enabled != nil {
		enabled = *request.Enabled
	}
	if request.Secret != nil {
		secret = strings.TrimSpace(*request.Secret)
	}
	if enabled && secret == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "启用服务间鉴权前必须配置密钥",
		})
		return
	}

	// Persist both values together so the database always contains a complete
	// relay authentication configuration.
	if err := model.UpdateOptionsBulk(map[string]string{
		common.RelayServiceAuthEnabledKey: boolToString(enabled),
		common.RelayServiceAuthSecretKey:  secret,
	}); err != nil {
		common.ApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    currentRelayServiceAuthResponse(),
	})
}
