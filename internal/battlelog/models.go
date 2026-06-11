package battlelog

import (
	"time"

	"trama/internal/core"
)

type PlayerID string

type Player struct {
	ID        PlayerID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TournamentID string
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

type Round struct {
	ID           string
	TournamentID TournamentID
	Number       RoundNumber
	CreatedAt    time.Time
}

type PairingID string
type PairingStatus string

const (
	PairingStatusPending PairingStatus = "pending"
	PairingStatusPlaying PairingStatus = "playing"
	PairingStatusDone    PairingStatus = "done"
)

type Pairing struct {
	ID        PairingID
	RoundID   string
	Player1ID PlayerID
	Player2ID PlayerID
	Status    PairingStatus
	CreatedAt time.Time
}

type BattleResultID string

type BattleResult struct {
	ID        BattleResultID
	PairingID PairingID
	PlayerID  PlayerID
	Score     int
	Winner    bool
	CreatedAt time.Time
}
