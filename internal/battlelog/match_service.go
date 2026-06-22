package battlelog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"trama/internal/user"
)

type matchRepository interface {
	Create(ctx context.Context, m Match) (Match, error)
	Get(ctx context.Context, id MatchID) (Match, error)
	List(ctx context.Context) ([]Match, error)
	Update(ctx context.Context, m Match) error
	Delete(ctx context.Context, id MatchID) error
}

type MatchService struct {
	repo matchRepository
}

func NewMatchService(repo matchRepository) *MatchService {
	return &MatchService{repo: repo}
}

func (s *MatchService) Create(ctx context.Context, in CreateMatchInput) (MatchOutput, error) {
	if in.Points1 < 0 || in.Points1 > 100 || in.Points2 < 0 || in.Points2 > 100 {
		return MatchOutput{}, ErrInvalidPoints
	}

	if in.WinnerID != in.Player1ID && in.WinnerID != in.Player2ID {
		return MatchOutput{}, ErrInvalidWinner
	}

	m := Match{
		Player1ID: user.NewUserID(in.Player1ID),
		Player2ID: user.NewUserID(in.Player2ID),
		Points1:   in.Points1,
		Points2:   in.Points2,
		WinnerID:  user.NewUserID(in.WinnerID),
	}

	created, err := s.repo.Create(ctx, m)
	if err != nil {
		return MatchOutput{}, err
	}

	return toMatchOutput(created), nil
}

func (s *MatchService) Get(ctx context.Context, id uuid.UUID) (MatchOutput, error) {
	m, err := s.repo.Get(ctx, MatchID{id})
	if err != nil {
		if errors.Is(err, ErrMatchNotFound) {
			return MatchOutput{}, fmt.Errorf("match %s: %w", id, ErrMatchNotFound)
		}
		return MatchOutput{}, err
	}

	return toMatchOutput(m), nil
}

func (s *MatchService) List(ctx context.Context) ([]MatchOutput, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]MatchOutput, len(items))
	for i, m := range items {
		out[i] = toMatchOutput(m)
	}
	return out, nil
}

func (s *MatchService) Update(ctx context.Context, in UpdateMatchInput) error {
	if in.Points1 < 0 || in.Points1 > 100 || in.Points2 < 0 || in.Points2 > 100 {
		return ErrInvalidPoints
	}

	if in.WinnerID != in.Player1ID && in.WinnerID != in.Player2ID {
		return ErrInvalidWinner
	}

	m := Match{
		ID:        MatchID{in.ID},
		Player1ID: user.NewUserID(in.Player1ID),
		Player2ID: user.NewUserID(in.Player2ID),
		Points1:   in.Points1,
		Points2:   in.Points2,
		WinnerID:  user.NewUserID(in.WinnerID),
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(ctx, m)
}

func (s *MatchService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, MatchID{id})
}

func toMatchOutput(m Match) MatchOutput {
	return MatchOutput{
		ID:        m.ID.UUID,
		Player1ID: m.Player1ID.UUID,
		Player2ID: m.Player2ID.UUID,
		Points1:   m.Points1,
		Points2:   m.Points2,
		WinnerID:  m.WinnerID.UUID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
