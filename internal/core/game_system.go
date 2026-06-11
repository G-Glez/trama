package core

import "time"

type GameSystemID string

type GameSystem struct {
	ID        GameSystemID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EditionID string

type Edition struct {
	ID           EditionID
	GameSystemID GameSystemID
	Name         string
	Version      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FactionID string

type Faction struct {
	ID        FactionID
	EditionID EditionID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
