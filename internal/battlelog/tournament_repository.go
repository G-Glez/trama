package battlelog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	battleloggen "trama/internal/gen/battlelog"
	"trama/internal/user"
)

type tournamentQuerier interface {
	CreateTournament(ctx context.Context, arg battleloggen.CreateTournamentParams) error
	GetTournament(ctx context.Context, id string) (battleloggen.Tournament, error)
	ListTournaments(ctx context.Context) ([]battleloggen.Tournament, error)
	UpdateTournament(ctx context.Context, arg battleloggen.UpdateTournamentParams) error
	DeleteTournament(ctx context.Context, id string) error
	AddTournamentPlayer(ctx context.Context, arg battleloggen.AddTournamentPlayerParams) error
	RemoveTournamentPlayer(ctx context.Context, arg battleloggen.RemoveTournamentPlayerParams) error
	ListTournamentPlayers(ctx context.Context, tournamentID string) ([]battleloggen.TournamentPlayer, error)
	DeleteTournamentPlayersByTournament(ctx context.Context, tournamentID string) error
}

type TournamentSQLRepository struct {
	querier tournamentQuerier
}

func NewTournamentRepository(querier tournamentQuerier) *TournamentSQLRepository {
	return &TournamentSQLRepository{querier: querier}
}

func (r *TournamentSQLRepository) Create(ctx context.Context, t Tournament) (Tournament, error) {
	t.ID = TournamentID{uuid.New()}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	playerID := sql.NullString{Valid: false}
	if t.PlayerID != nil {
		playerID = sql.NullString{String: t.PlayerID.String(), Valid: true}
	}

	err := r.querier.CreateTournament(ctx, battleloggen.CreateTournamentParams{
		ID:        t.ID.String(),
		Name:      t.Name,
		PlayerID:  playerID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	})
	if err != nil {
		return Tournament{}, fmt.Errorf("database error: %w", err)
	}

	return t, nil
}

func (r *TournamentSQLRepository) Get(ctx context.Context, id TournamentID) (Tournament, error) {
	t, err := r.querier.GetTournament(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Tournament{}, ErrTournamentNotFound
		}
		return Tournament{}, fmt.Errorf("database error: %w", err)
	}

	return toDomainTournament(t), nil
}

func (r *TournamentSQLRepository) List(ctx context.Context) ([]Tournament, error) {
	items, err := r.querier.ListTournaments(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	list := make([]Tournament, len(items))
	for i, t := range items {
		list[i] = toDomainTournament(t)
	}
	return list, nil
}

func (r *TournamentSQLRepository) Update(ctx context.Context, t Tournament) error {
	playerID := sql.NullString{Valid: false}
	if t.PlayerID != nil {
		playerID = sql.NullString{String: t.PlayerID.String(), Valid: true}
	}

	err := r.querier.UpdateTournament(ctx, battleloggen.UpdateTournamentParams{
		Name:      t.Name,
		PlayerID:  playerID,
		UpdatedAt: t.UpdatedAt,
		ID:        t.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *TournamentSQLRepository) Delete(ctx context.Context, id TournamentID) error {
	err := r.querier.DeleteTournament(ctx, id.String())
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *TournamentSQLRepository) AddPlayer(ctx context.Context, tp TournamentPlayer) error {
	teamID := sql.NullString{Valid: false}
	if tp.TeamID != nil {
		teamID = sql.NullString{String: tp.TeamID.String(), Valid: true}
	}

	err := r.querier.AddTournamentPlayer(ctx, battleloggen.AddTournamentPlayerParams{
		TournamentID: tp.TournamentID.String(),
		PlayerID:     tp.PlayerID.String(),
		TeamID:       teamID,
	})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *TournamentSQLRepository) RemovePlayer(ctx context.Context, tournamentID TournamentID, playerID user.UserID) error {
	err := r.querier.RemoveTournamentPlayer(ctx, battleloggen.RemoveTournamentPlayerParams{
		TournamentID: tournamentID.String(),
		PlayerID:     playerID.String(),
	})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *TournamentSQLRepository) ListPlayers(ctx context.Context, tournamentID TournamentID) ([]TournamentPlayer, error) {
	items, err := r.querier.ListTournamentPlayers(ctx, tournamentID.String())
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	list := make([]TournamentPlayer, len(items))
	for i, tp := range items {
		list[i] = toDomainTournamentPlayer(tp)
	}
	return list, nil
}

func (r *TournamentSQLRepository) DeletePlayersByTournament(ctx context.Context, tournamentID TournamentID) error {
	err := r.querier.DeleteTournamentPlayersByTournament(ctx, tournamentID.String())
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func toDomainTournament(t battleloggen.Tournament) Tournament {
	var playerID *user.UserID
	if t.PlayerID.Valid {
		uid := user.NewUserID(uuid.MustParse(t.PlayerID.String))
		playerID = &uid
	}

	return Tournament{
		ID:        TournamentID{uuid.MustParse(t.ID)},
		Name:      t.Name,
		PlayerID:  playerID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func toDomainTournamentPlayer(tp battleloggen.TournamentPlayer) TournamentPlayer {
	var teamID *TeamID
	if tp.TeamID.Valid {
		tid := TeamID{uuid.MustParse(tp.TeamID.String)}
		teamID = &tid
	}

	return TournamentPlayer{
		TournamentID: TournamentID{uuid.MustParse(tp.TournamentID)},
		PlayerID:     user.NewUserID(uuid.MustParse(tp.PlayerID)),
		TeamID:       teamID,
	}
}
