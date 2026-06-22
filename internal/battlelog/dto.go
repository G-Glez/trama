package battlelog

import (
	"time"

	"github.com/google/uuid"
)

type CreateMatchInput struct {
	Player1ID uuid.UUID `json:"player1_id" binding:"required"`
	Player2ID uuid.UUID `json:"player2_id" binding:"required"`
	Points1   int       `json:"points1"    binding:"required,min=0,max=100"`
	Points2   int       `json:"points2"    binding:"required,min=0,max=100"`
	WinnerID  uuid.UUID `json:"winner_id"  binding:"required"`
}

type UpdateMatchInput struct {
	ID        uuid.UUID `json:"-"`
	Player1ID uuid.UUID `json:"player1_id" binding:"required"`
	Player2ID uuid.UUID `json:"player2_id" binding:"required"`
	Points1   int       `json:"points1"    binding:"required,min=0,max=100"`
	Points2   int       `json:"points2"    binding:"required,min=0,max=100"`
	WinnerID  uuid.UUID `json:"winner_id"  binding:"required"`
}

type MatchOutput struct {
	ID        uuid.UUID `json:"id"`
	Player1ID uuid.UUID `json:"player1_id"`
	Player2ID uuid.UUID `json:"player2_id"`
	Points1   int       `json:"points1"`
	Points2   int       `json:"points2"`
	WinnerID  uuid.UUID `json:"winner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTeamInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTeamInput struct {
	ID   uuid.UUID `json:"-"`
	Name string    `json:"name" binding:"required"`
}

type TeamOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTournamentInput struct {
	Name     string     `json:"name"      binding:"required"`
	PlayerID *uuid.UUID `json:"player_id"`
}

type UpdateTournamentInput struct {
	ID       uuid.UUID  `json:"-"`
	Name     string     `json:"name"      binding:"required"`
	PlayerID *uuid.UUID `json:"player_id"`
}

type TournamentOutput struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	PlayerID  *uuid.UUID `json:"player_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
