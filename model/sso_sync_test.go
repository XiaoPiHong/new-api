package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSSOSyncRequiresMatchingOrdinaryOIDCOwner(t *testing.T) {
	owner := User{Id: 42, OidcId: "7", Role: common.RoleCommonUser, Status: common.UserStatusEnabled, AuthVersion: 9}
	input := SSOUserSyncInput{Subject: "7", ExistingUserID: 42}
	require.NoError(t, validateSSOSyncIdentity(&owner, input))
	for _, tc := range []struct {
		name   string
		change func(*User, *SSOUserSyncInput)
	}{
		{"unbound same-name account", func(user *User, _ *SSOUserSyncInput) { user.OidcId = "" }},
		{"another subject", func(user *User, _ *SSOUserSyncInput) { user.OidcId = "8" }},
		{"stale mapping", func(_ *User, request *SSOUserSyncInput) { request.ExistingUserID = 99 }},
		{"administrator", func(user *User, _ *SSOUserSyncInput) { user.Role = common.RoleAdminUser }},
		{"root", func(user *User, _ *SSOUserSyncInput) { user.Role = common.RoleRootUser }},
		{"disabled", func(user *User, _ *SSOUserSyncInput) { user.Status = common.UserStatusDisabled }},
		{"deleted", func(user *User, _ *SSOUserSyncInput) { user.DeletedAt.Valid = true }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			user, request := owner, input
			tc.change(&user, &request)
			assert.ErrorIs(t, validateSSOSyncIdentity(&user, request), ErrSSOSyncIdentity)
			assert.Equal(t, int64(9), user.AuthVersion)
		})
	}
}

func TestSSOSyncRejectsAnElevatedOrUnboundedPATBeforeDatabaseAccess(t *testing.T) {
	input := SSOUserSyncInput{Subject: "7", Username: "sso-user", ActorID: 1}
	token := &UserAccessToken{TokenHash: "test-fingerprint", ExpiresAt: common.GetTimestamp() + 3600}
	require.NoError(t, token.SetScopes([]string{"account_security:write", "user:write"}))
	_, _, err := SyncSSODashboardUser(input, token, 20)
	assert.ErrorIs(t, err, ErrSSOSyncRequest)
	require.NoError(t, token.SetScopes(SSODashboardScopes()))
	token.ExpiresAt = 0
	_, _, err = SyncSSODashboardUser(input, token, 20)
	assert.ErrorIs(t, err, ErrSSOSyncRequest)
}

func TestSSOSyncImportsLegacyOwnerInitializesEmailAndRotatesOnlyItsPAT(t *testing.T) {
	truncateTables(t)
	root := User{Username: "sso-root", Password: "test-hash", AffCode: "sso-root", Role: common.RoleRootUser, Status: common.UserStatusEnabled}
	owner := User{Username: "legacy-owner", Password: "test-hash", AffCode: "sso-owner", OidcId: "7", Role: common.RoleCommonUser, Status: common.UserStatusEnabled, AuthVersion: 9}
	require.NoError(t, DB.Create(&root).Error)
	require.NoError(t, DB.Create(&owner).Error)
	manual := UserAccessToken{UserId: owner.Id, Name: "manual", TokenHash: AccessTokenFingerprint("manual")}
	require.NoError(t, manual.SetScopes([]string{"profile:read"}))
	require.NoError(t, DB.Create(&manual).Error)
	input := SSOUserSyncInput{ActorID: root.Id, Subject: "7", Username: owner.Username, ExistingUserID: owner.Id, Email: "first@example.com"}
	first := &UserAccessToken{TokenHash: AccessTokenFingerprint("first-sso"), ExpiresAt: common.GetTimestamp() + 3600}
	require.NoError(t, first.SetScopes(SSODashboardScopes()))
	user, created, err := SyncSSODashboardUser(input, first, 20)
	require.NoError(t, err)
	assert.False(t, created)
	assert.Equal(t, "first@example.com", user.Email)
	assert.Equal(t, int64(9), user.AuthVersion)
	var claim ExternalIdentityClaim
	require.NoError(t, DB.Where("provider = ? AND subject = ?", SSOIdentityProvider, "7").First(&claim).Error)
	assert.Equal(t, owner.Id, claim.UserId)

	input.PreviousTokenRef = first.TokenHash
	input.Email = "replacement@example.com"
	second := &UserAccessToken{TokenHash: AccessTokenFingerprint("second-sso"), ExpiresAt: common.GetTimestamp() + 3600}
	require.NoError(t, second.SetScopes(SSODashboardScopes()))
	user, _, err = SyncSSODashboardUser(input, second, 20)
	require.NoError(t, err)
	assert.Equal(t, "first@example.com", user.Email)
	assert.Equal(t, int64(9), user.AuthVersion)
	var tokens []UserAccessToken
	require.NoError(t, DB.Where("user_id = ?", owner.Id).Order("id").Find(&tokens).Error)
	require.Len(t, tokens, 2)
	assert.Equal(t, manual.TokenHash, tokens[0].TokenHash)
	assert.Equal(t, second.TokenHash, tokens[1].TokenHash)
}

func TestSSOSyncConflictingDurableClaimRollsBackNewUserAndPAT(t *testing.T) {
	truncateTables(t)
	root := User{Username: "sso-root", Password: "test-hash", AffCode: "sso-root", Role: common.RoleRootUser, Status: common.UserStatusEnabled}
	claimOwner := User{Username: "claim-owner", Password: "test-hash", AffCode: "claim-owner", OidcId: "8", Role: common.RoleCommonUser, Status: common.UserStatusEnabled}
	require.NoError(t, DB.Create(&root).Error)
	require.NoError(t, DB.Create(&claimOwner).Error)
	require.NoError(t, DB.Transaction(func(tx *gorm.DB) error {
		return ClaimExternalIdentityWithTx(tx, SSOIdentityProvider, "99", claimOwner.Id)
	}))
	token := &UserAccessToken{TokenHash: AccessTokenFingerprint("new-sso"), ExpiresAt: common.GetTimestamp() + 3600}
	require.NoError(t, token.SetScopes(SSODashboardScopes()))
	_, _, err := SyncSSODashboardUser(SSOUserSyncInput{ActorID: root.Id, Subject: "99", Username: "new-owner"}, token, 20)
	assert.ErrorIs(t, err, ErrExternalIdentityAlreadyClaimed)
	var count int64
	require.NoError(t, DB.Model(&User{}).Where("username = ?", "new-owner").Count(&count).Error)
	assert.Zero(t, count)
	require.NoError(t, DB.Model(&UserAccessToken{}).Count(&count).Error)
	assert.Zero(t, count)
}
