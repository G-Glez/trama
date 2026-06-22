package battlelog

import (
	"time"

	"github.com/google/uuid"

	"trama/internal/user"
)

type MatchID struct{ uuid.UUID }

type Match struct {
	ID        MatchID
	Player1ID user.UserID
	Player2ID user.UserID
	Points1   int
	Points2   int
	WinnerID  user.UserID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TeamID struct{ uuid.UUID }

type Team struct {
	ID        TeamID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TournamentID struct{ uuid.UUID }

type Tournament struct {
	ID        TournamentID
	Name      string
	PlayerID  *user.UserID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TournamentPlayer struct {
	TournamentID TournamentID
	PlayerID     user.UserID
	TeamID       *TeamID
}
