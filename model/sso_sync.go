package model

import (
	"errors"
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

var (
	ErrSSOSyncIdentity = errors.New("SSO identity does not own this account")
	ErrSSOSyncRequest  = errors.New("invalid SSO synchronization request")
)

const SSOIdentityProvider = "sso"

type SSOUserSyncInput struct {
	Subject          string `json:"subject"`
	Username         string `json:"username"`
	DisplayName      string `json:"display_name"`
	Email            string `json:"email"`
	ExistingUserID   int    `json:"existing_user_id,omitempty"`
	PreviousTokenRef string `json:"previous_token_ref,omitempty"`
	ActorID          int    `json:"-"`
}

// SSODashboardScopes excludes account-security and administrative grants.
func SSODashboardScopes() []string {
	return []string{
		"api_key:read", "api_key:reveal", "api_key:write",
		"profile:read", "profile:write", "usage:read", "wallet:read", "wallet:write",
	}
}

func validateSSOSyncIdentity(user *User, input SSOUserSyncInput) error {
	if user == nil || user.DeletedAt.Valid || user.Role != common.RoleCommonUser ||
		user.Status != common.UserStatusEnabled || user.OidcId != input.Subject ||
		(input.ExistingUserID > 0 && user.Id != input.ExistingUserID) {
		return ErrSSOSyncIdentity
	}
	return nil
}

// SyncSSODashboardUser never claims an existing account by username or email.
// Creation, ownership checks, and scoped PAT rotation share one transaction.
func SyncSSODashboardUser(input SSOUserSyncInput, token *UserAccessToken, limit int) (*User, bool, error) {
	subjectID, err := strconv.ParseUint(input.Subject, 10, 64)
	now := common.GetTimestamp()
	if err != nil || subjectID == 0 || strconv.FormatUint(subjectID, 10) != input.Subject ||
		input.ActorID <= 0 || input.ExistingUserID < 0 || input.Username == "" ||
		utf8.RuneCountInString(input.Username) > 20 || utf8.RuneCountInString(input.DisplayName) > 20 ||
		utf8.RuneCountInString(input.Email) > 50 || token == nil || token.TokenHash == "" ||
		token.ExpiresAt <= now || token.ExpiresAt > now+31*24*60*60 || limit <= 0 ||
		!slices.Equal(token.GetScopes(), SSODashboardScopes()) {
		return nil, false, ErrSSOSyncRequest
	}

	var user User
	created := false
	emailInitialized := false
	err = DB.Transaction(func(tx *gorm.DB) error {
		// Match the official email-binding lock order before locking users.
		if err := lockNormalizedEmail(tx, input.Email); err != nil {
			return err
		}
		// Serialize PAT changes by this issuer. The durable subject claim below
		// also prevents competing issuers from creating different owners.
		var actor User
		if err := lockForUpdate(tx).Select("id", "role", "status").First(&actor, input.ActorID).Error; err != nil {
			return err
		}
		if actor.Role != common.RoleRootUser || actor.Status != common.UserStatusEnabled {
			return ErrSSOSyncIdentity
		}
		var matches []User
		if err := lockForUpdate(tx.Unscoped()).Where("oidc_id = ?", input.Subject).Limit(2).Find(&matches).Error; err != nil {
			return err
		}
		if len(matches) > 1 {
			return ErrSSOSyncIdentity
		}
		if len(matches) == 1 {
			user = matches[0]
			if err := validateSSOSyncIdentity(&user, input); err != nil {
				return err
			}
		} else {
			if input.ExistingUserID > 0 {
				return ErrSSOSyncIdentity
			}
			var usernameCount int64
			if err := tx.Unscoped().Model(&User{}).Where("username = ?", input.Username).Count(&usernameCount).Error; err != nil {
				return err
			}
			if usernameCount > 0 {
				return ErrSSOSyncIdentity
			}
			password, err := common.GenerateRandomCharsKey(43)
			if err != nil {
				return err
			}
			user = User{
				Username: input.Username, DisplayName: input.DisplayName, Password: password,
				Email: input.Email, OidcId: input.Subject, Group: "default",
				Role: common.RoleCommonUser, Status: common.UserStatusEnabled, AuthVersion: 1,
			}
			if err := user.InsertWithTx(tx, 0); err != nil {
				return err
			}
			created = true
		}

		// Import an exact legacy OIDC mapping idempotently; never overwrite a
		// claim held by another user or another subject in this provider slot.
		if err := ClaimExternalIdentityWithTx(tx, SSOIdentityProvider, input.Subject, user.Id); err != nil {
			return err
		}
		if !created && user.Email == "" && input.Email != "" {
			if err := ensureEmailAvailableWithTx(tx, input.Email, user.Id); err != nil {
				return err
			}
			result := tx.Model(&User{}).Where("id = ? AND (email = ? OR email IS NULL)", user.Id, "").Update("email", input.Email)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 1 {
				user.Email = input.Email
				emailInitialized = true
			}
		}
		tokenName := "SSO dashboard/" + input.Subject
		if input.PreviousTokenRef != "" {
			// A stale cache must never revoke another application's PAT.
			if err := tx.Where("user_id = ? AND token_hash = ? AND name = ? AND scopes = ?",
				user.Id, input.PreviousTokenRef, tokenName, token.Scopes).Delete(&UserAccessToken{}).Error; err != nil {
				return err
			}
		}
		var count int64
		if err := tx.Model(&UserAccessToken{}).Where("user_id = ?", user.Id).Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(limit) {
			return ErrAccessTokenLimit
		}
		token.Id = 0
		token.UserId = user.Id
		token.Name = tokenName
		token.CreatedAt = now
		return tx.Create(token).Error
	})
	if err != nil {
		return nil, false, err
	}
	if created {
		user.FinishInsert(0)
	} else if emailInitialized {
		if err := updateUserCache(user); err != nil {
			return nil, false, err
		}
	}
	return &user, created, nil
}
