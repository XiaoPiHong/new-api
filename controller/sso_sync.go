package controller

import (
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

// SyncSSOUser is a root-authorized service operation, not a user token-generation
// or account-linking endpoint. It also requires the configured service secret.
func SyncSSOUser(c *gin.Context) {
	expected := strings.TrimSpace(common.GetEnvOrDefaultString("BILLING_TRACE_SECRET", ""))
	actual := strings.TrimSpace(c.GetHeader(common.BillingTraceSecretKey))
	if c.GetInt("role") != common.RoleRootUser || expected == "" ||
		len(expected) != len(actual) || subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) != 1 {
		common.ApiErrorStatus(c, http.StatusForbidden, common.NewMessage("SSO synchronization is not authorized."), gin.H{"code": "SSO_SYNC_FORBIDDEN"})
		return
	}
	var input model.SSOUserSyncInput
	if err := common.DecodeJson(c.Request.Body, &input); err != nil {
		common.ApiErrorStatus(c, http.StatusBadRequest, common.NewMessage("Invalid SSO synchronization request."), gin.H{"code": "SSO_SYNC_INVALID_REQUEST"})
		return
	}
	input.Subject = strings.TrimSpace(input.Subject)
	input.Username = strings.TrimSpace(input.Username)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Email = model.NormalizeEmail(input.Email)
	input.PreviousTokenRef = strings.ToLower(strings.TrimSpace(input.PreviousTokenRef))
	input.ActorID = c.GetInt("id")
	if input.DisplayName == "" {
		input.DisplayName = input.Username
	}
	if input.Email != "" {
		address, err := mail.ParseAddress(input.Email)
		if err != nil || address.Address != input.Email {
			common.ApiErrorStatus(c, http.StatusBadRequest, common.NewMessage("Invalid SSO synchronization request."), gin.H{"code": "SSO_SYNC_INVALID_REQUEST"})
			return
		}
	}
	if input.PreviousTokenRef != "" {
		decoded, err := hex.DecodeString(input.PreviousTokenRef)
		if err != nil || len(decoded) != 32 {
			common.ApiErrorStatus(c, http.StatusBadRequest, common.NewMessage("Invalid SSO synchronization request."), gin.H{"code": "SSO_SYNC_INVALID_REQUEST"})
			return
		}
	}
	scopes, err := service.NormalizeAccessTokenScopes(0, common.RoleCommonUser, model.SSODashboardScopes())
	if err != nil {
		writeSecurityOperationError(c, err)
		return
	}
	suffix, err := common.GenerateRandomCharsKey(43)
	if err != nil {
		writeSecurityOperationError(c, err)
		return
	}
	raw := model.AccessTokenPrefix + suffix
	token := &model.UserAccessToken{
		TokenHash: model.AccessTokenFingerprint(raw), TokenHint: model.AccessTokenHint(raw),
		ExpiresAt: common.GetTimestamp() + service.AccessTokenDefaultExpiryDays*24*60*60,
	}
	if err := token.SetScopes(scopes); err != nil {
		writeSecurityOperationError(c, err)
		return
	}
	user, created, err := model.SyncSSODashboardUser(input, token, service.AccessTokenMaxPerUser)
	switch {
	case errors.Is(err, model.ErrSSOSyncIdentity), errors.Is(err, model.ErrExternalIdentityAlreadyClaimed), errors.Is(err, model.ErrEmailAlreadyTaken):
		common.ApiErrorStatus(c, http.StatusConflict, common.NewMessage("SSO identity does not own this account."), gin.H{"code": "SSO_SYNC_IDENTITY_CONFLICT"})
		return
	case errors.Is(err, model.ErrSSOSyncRequest):
		common.ApiErrorStatus(c, http.StatusBadRequest, common.NewMessage("Invalid SSO synchronization request."), gin.H{"code": "SSO_SYNC_INVALID_REQUEST"})
		return
	case errors.Is(err, model.ErrAccessTokenLimit):
		common.ApiErrorStatus(c, http.StatusConflict, common.NewMessage("The user has reached the access token limit."), gin.H{"code": "SSO_SYNC_TOKEN_LIMIT"})
		return
	case err != nil:
		common.ApiErrorStatus(c, http.StatusInternalServerError, common.NewMessage("SSO synchronization failed."), gin.H{"code": "SSO_SYNC_FAILED"})
		return
	}
	recordManageAuditFor(c, user.Id, "access_token.generate", map[string]any{
		"token_id": token.Id, "name": token.Name, "scopes": scopes,
		"expires_at": token.ExpiresAt, "token_ref": token.TokenHash,
	})
	common.ApiSuccess(c, gin.H{
		"user_id": user.Id, "username": user.Username, "access_token": raw,
		"expires_at": token.ExpiresAt, "created": created,
	})
}
