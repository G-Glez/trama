package core

import "time"

type GameSystemOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditionOutput struct {
	ID           string    `json:"id"`
	GameSystemID string    `json:"game_system_id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FactionOutput struct {
	ID        string    `json:"id"`
	EditionID string    `json:"edition_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateGameSystemInput struct {
	Name string `json:"name"`
}

type UpdateGameSystemInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateEditionInput struct {
	GameSystemID string `json:"game_system_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
}

type UpdateEditionInput struct {
	ID           string `json:"id"`
	GameSystemID string `json:"game_system_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
}

type CreateFactionInput struct {
	EditionID string `json:"edition_id"`
	Name      string `json:"name"`
}

type UpdateFactionInput struct {
	ID        string `json:"id"`
	EditionID string `json:"edition_id"`
	Name      string `json:"name"`
}
