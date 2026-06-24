package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// -----------------------------------------------------------------------------------
// JWT implements TokenProvider using HS256 HMAC signing.
// -----------------------------------------------------------------------------------
type JWT struct {
	secret []byte
}

func NewJWT(secret []byte) *JWT {
	return &JWT{secret: secret}
}

// -----------------------------------------------------------------------------------
// Generate signs the given claims into a JWT string using HS256.
// Input: claims to be signed
// Output: signed JWT string
// -----------------------------------------------------------------------------------
func (j *JWT) Generate(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signed, nil
}

// -----------------------------------------------------------------------------------
// Validate parses and verifies a JWT string.
// Input: JWT string to validate
// Output: Claims
// -----------------------------------------------------------------------------------
func (j *JWT) Validate(tokenString string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("parse jwt: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}

	return *claims, nil
}
