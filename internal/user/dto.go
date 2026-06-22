package user

import (
	"time"

	"github.com/google/uuid"
)

type RegisterInput struct {
	Username       string    `json:"username" binding:"required"`
	Email          string    `json:"email"    binding:"required,email"`
	Password       string    `json:"password" binding:"required,min=8"`
	DefaultFaction uuid.UUID `json:"default_faction"`
}

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	ID             uuid.UUID `json:"-"`
	Username       string    `json:"username"        binding:"required"`
	Email          string    `json:"email"           binding:"required,email"`
	Password       string    `json:"password,omitempty"`
	DefaultFaction uuid.UUID `json:"default_faction"`
}

type UserOutput struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	DefaultFaction uuid.UUID `json:"default_faction"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
