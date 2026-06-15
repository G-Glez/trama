package core

import (
	"time"

	"github.com/google/uuid"
)

type GameSystemOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditionOutput struct {
	ID           uuid.UUID `json:"id"`
	GameSystemID uuid.UUID `json:"game_system_id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FactionOutput struct {
	ID        uuid.UUID `json:"id"`
	EditionID uuid.UUID `json:"edition_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateGameSystemInput struct {
	Name string `json:"name"`
}

type UpdateGameSystemInput struct {
	ID   uuid.UUID `json:"-"`
	Name string    `json:"name"`
}

type CreateEditionInput struct {
	GameSystemID uuid.UUID `json:"game_system_id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
}

type UpdateEditionInput struct {
	ID           uuid.UUID `json:"-"`
	GameSystemID uuid.UUID `json:"game_system_id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
}

type CreateFactionInput struct {
	EditionID uuid.UUID `json:"edition_id"`
	Name      string    `json:"name"`
}

type UpdateFactionInput struct {
	ID        uuid.UUID `json:"-"`
	EditionID uuid.UUID `json:"edition_id"`
	Name      string    `json:"name"`
}
