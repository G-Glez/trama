package core

import (
	"time"

	"github.com/google/uuid"
)

type GameSystemID struct {
	uuid.UUID
}

type GameSystem struct {
	ID        GameSystemID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EditionID struct {
	uuid.UUID
}

type Edition struct {
	ID           EditionID
	GameSystemID GameSystemID
	Name         string
	Version      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FactionID struct {
	uuid.UUID
}

type Faction struct {
	ID        FactionID
	EditionID EditionID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
