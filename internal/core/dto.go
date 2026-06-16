package core

import (
	"time"

	"github.com/google/uuid"
)

type FactionOutput struct {
	ID           uuid.UUID `json:"id"`
	EditionID    uuid.UUID `json:"edition_id"`
	GameSystemID uuid.UUID `json:"game_system_id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateFactionInput struct {
	EditionID    uuid.UUID `json:"edition_id"`
	GameSystemID uuid.UUID `json:"game_system_id"`
	Name         string    `json:"name"`
}

type UpdateFactionInput struct {
	ID           uuid.UUID `json:"-"`
	EditionID    uuid.UUID `json:"edition_id"`
	GameSystemID uuid.UUID `json:"game_system_id"`
	Name         string    `json:"name"`
}
