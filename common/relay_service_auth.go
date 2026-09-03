package common

import "strings"

const (
	RelayServiceAuthEnabledKey = "relay_service_auth.enabled"
	RelayServiceAuthSecretKey  = "relay_service_auth.secret"
)

// GetRelayServiceAuthConfig returns the current relay protection state.
// The values are loaded from New API's database-backed OptionMap. The secret
// is intentionally kept server-side.
func GetRelayServiceAuthConfig() (enabled bool, secret string) {
	OptionMapRWMutex.RLock()
	if value, ok := OptionMap[RelayServiceAuthEnabledKey]; ok {
		normalized := strings.TrimSpace(value)
		enabled = strings.EqualFold(normalized, "true") || normalized == "1"
	}
	if value, ok := OptionMap[RelayServiceAuthSecretKey]; ok {
		secret = strings.TrimSpace(value)
	}
	OptionMapRWMutex.RUnlock()
	return enabled, secret
}
