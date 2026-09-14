package auth

import (
	"crypto/subtle"
)

// APIKeyValidator handles validating static API keys
type APIKeyValidator struct {
	validKeys map[string]bool
}

// NewAPIKeyValidator creates a new validator with the given valid keys
func NewAPIKeyValidator(keys []string) *APIKeyValidator {
	validKeys := make(map[string]bool)
	for _, k := range keys {
		validKeys[k] = true
	}
	return &APIKeyValidator{
		validKeys: validKeys,
	}
}

// ValidateKey validates an API key in constant time (if single key) or maps
func (v *APIKeyValidator) ValidateKey(key string) bool {
	if key == "" {
		return false
	}
	
	// Fast path for map lookup, but vulnerable to timing attacks if we cared
	// For better security, we could do constant-time comparisons against all keys
	for validKey := range v.validKeys {
		if subtle.ConstantTimeCompare([]byte(key), []byte(validKey)) == 1 {
			return true
		}
	}
	
	return false
}
