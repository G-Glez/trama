package battlelog

import (
	"time"

	"github.com/google/uuid"

	"trama/internal/core"
)

type PlayerID struct {
	uuid.UUID
}

type Player struct {
	ID        PlayerID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TournamentID struct {
	uuid.UUID
}
type TournamentStatus string

const (
	TournamentStatusPlanned  TournamentStatus = "planned"
	TournamentStatusOngoing  TournamentStatus = "ongoing"
	TournamentStatusFinished TournamentStatus = "finished"
)

type Tournament struct {
	ID         TournamentID
	Name       string
	GameSystem core.GameSystemID
	Edition    core.EditionID
	Status     TournamentStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type RoundNumber int

type RoundID struct {
	uuid.UUID
}

type Round struct {
	ID           RoundID
	TournamentID TournamentID
	Number       RoundNumber
	CreatedAt    time.Time
}

type PairingID struct {
	uuid.UUID
}
type PairingStatus string

const (
	PairingStatusPending PairingStatus = "pending"
	PairingStatusPlaying PairingStatus = "playing"
	PairingStatusDone    PairingStatus = "done"
)

type Pairing struct {
	ID        PairingID
	RoundID   RoundID
	Player1ID PlayerID
	Player2ID PlayerID
	Status    PairingStatus
	CreatedAt time.Time
}

type BattleResultID struct {
	uuid.UUID
}

type BattleResult struct {
	ID        BattleResultID
	PairingID PairingID
	PlayerID  PlayerID
	Score     int
	Winner    bool
	CreatedAt time.Time
}
