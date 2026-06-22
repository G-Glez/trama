package user

import (
	"time"

	"github.com/google/uuid"
)

type UserID struct{ uuid.UUID }

func NewUserID(id uuid.UUID) UserID { return UserID{id} }

type User struct {
	ID             UserID
	Username       string
	Email          string
	PasswordHash   string
	DefaultFaction uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
