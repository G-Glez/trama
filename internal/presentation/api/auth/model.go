package auth

import "github.com/golang-jwt/jwt/v5"

// -----------------------------------------------------------------------------------
// Claims carries the JWT payload for an authenticated user.
// -----------------------------------------------------------------------------------
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// -----------------------------------------------------------------------------------
// User represents a persisted user record in DynamoDB. Username is the primary key.
// -----------------------------------------------------------------------------------
type User struct {
	Username     string
	PasswordHash string
	CreatedAt    string
}
