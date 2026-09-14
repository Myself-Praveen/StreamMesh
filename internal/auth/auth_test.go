package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTValidator(t *testing.T) {
	secret := "testsecret"
	validator := NewJWTValidator(secret)

	// Create valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	claims, err := validator.ValidateToken(tokenString)
	if err != nil {
		t.Errorf("Expected valid token, got error: %v", err)
	}
	if claims["sub"] != "user123" {
		t.Errorf("Expected sub=user123, got %v", claims["sub"])
	}

	// Test invalid token
	_, err = validator.ValidateToken("invalid.token.string")
	if err == nil {
		t.Errorf("Expected error for invalid token")
	}
}

func TestAPIKeyValidator(t *testing.T) {
	validator := NewAPIKeyValidator([]string{"key1", "key2"})

	if !validator.ValidateKey("key1") {
		t.Errorf("Expected key1 to be valid")
	}
	if !validator.ValidateKey("key2") {
		t.Errorf("Expected key2 to be valid")
	}
	if validator.ValidateKey("key3") {
		t.Errorf("Expected key3 to be invalid")
	}
	if validator.ValidateKey("") {
		t.Errorf("Expected empty key to be invalid")
	}
}
