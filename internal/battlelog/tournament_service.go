package battlelog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"trama/internal/user"
)

type tournamentRepository interface {
	Create(ctx context.Context, t Tournament) (Tournament, error)
	Get(ctx context.Context, id TournamentID) (Tournament, error)
	List(ctx context.Context) ([]Tournament, error)
	Update(ctx context.Context, t Tournament) error
	Delete(ctx context.Context, id TournamentID) error
	AddPlayer(ctx context.Context, tp TournamentPlayer) error
	RemovePlayer(ctx context.Context, tournamentID TournamentID, playerID user.UserID) error
	ListPlayers(ctx context.Context, tournamentID TournamentID) ([]TournamentPlayer, error)
	DeletePlayersByTournament(ctx context.Context, tournamentID TournamentID) error
}

type TournamentService struct {
	repo tournamentRepository
}

func NewTournamentService(repo tournamentRepository) *TournamentService {
	return &TournamentService{repo: repo}
}

func (s *TournamentService) Create(ctx context.Context, in CreateTournamentInput) (TournamentOutput, error) {
	var playerID *user.UserID
	if in.PlayerID != nil {
		uid := user.NewUserID(*in.PlayerID)
		playerID = &uid
	}

	t := Tournament{
		Name:     in.Name,
		PlayerID: playerID,
	}

	created, err := s.repo.Create(ctx, t)
	if err != nil {
		return TournamentOutput{}, err
	}

	return toTournamentOutput(created), nil
}

func (s *TournamentService) Get(ctx context.Context, id uuid.UUID) (TournamentOutput, error) {
	t, err := s.repo.Get(ctx, TournamentID{id})
	if err != nil {
		if errors.Is(err, ErrTournamentNotFound) {
			return TournamentOutput{}, fmt.Errorf("tournament %s: %w", id, ErrTournamentNotFound)
		}
		return TournamentOutput{}, err
	}

	return toTournamentOutput(t), nil
}

func (s *TournamentService) List(ctx context.Context) ([]TournamentOutput, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]TournamentOutput, len(items))
	for i, t := range items {
		out[i] = toTournamentOutput(t)
	}
	return out, nil
}

func (s *TournamentService) Update(ctx context.Context, in UpdateTournamentInput) error {
	var playerID *user.UserID
	if in.PlayerID != nil {
		uid := user.NewUserID(*in.PlayerID)
		playerID = &uid
	}

	t := Tournament{
		ID:        TournamentID{in.ID},
		Name:      in.Name,
		PlayerID:  playerID,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(ctx, t)
}

func (s *TournamentService) Delete(ctx context.Context, id uuid.UUID) error {
	tID := TournamentID{id}

	if err := s.repo.DeletePlayersByTournament(ctx, tID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, tID)
}

func toTournamentOutput(t Tournament) TournamentOutput {
	var pid *uuid.UUID
	if t.PlayerID != nil {
		uid := t.PlayerID.UUID
		pid = &uid
	}

	return TournamentOutput{
		ID:        t.ID.UUID,
		Name:      t.Name,
		PlayerID:  pid,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
