package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWTValidator handles validating JWT tokens
type JWTValidator struct {
	secretKey []byte
}

// NewJWTValidator creates a new validator with the given secret
func NewJWTValidator(secret string) *JWTValidator {
	return &JWTValidator{
		secretKey: []byte(secret),
	}
}

// ValidateToken validates a JWT string and returns the claims
func (v *JWTValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
