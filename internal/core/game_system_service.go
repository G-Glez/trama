package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type gameSystemRepository interface {
	Create(ctx context.Context, gs GameSystem) (GameSystem, error)
	Get(ctx context.Context, id GameSystemID) (GameSystem, error)
	GetAll(ctx context.Context) ([]GameSystem, error)
	Update(ctx context.Context, gs GameSystem) error
	Delete(ctx context.Context, id GameSystemID) error
}

type GameSystemService struct {
	repo gameSystemRepository
}

func NewGameSystemService(repo gameSystemRepository) *GameSystemService {
	return &GameSystemService{repo: repo}
}

func (s *GameSystemService) Create(ctx context.Context, in CreateGameSystemInput) (GameSystemOutput, error) {
	gs := GameSystem{Name: in.Name}
	created, err := s.repo.Create(ctx, gs)
	if err != nil {
		return GameSystemOutput{}, err
	}

	return toGameSystemOutput(created), nil
}

func (s *GameSystemService) Get(ctx context.Context, id uuid.UUID) (GameSystemOutput, error) {
	gs, err := s.repo.Get(ctx, GameSystemID{id})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return GameSystemOutput{}, fmt.Errorf("game system %s: %w", id, ErrNotFound)
		}
		return GameSystemOutput{}, err
	}

	return toGameSystemOutput(gs), nil
}

func (s *GameSystemService) GetAll(ctx context.Context) ([]GameSystemOutput, error) {
	items, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]GameSystemOutput, len(items))
	for i, gs := range items {
		out[i] = toGameSystemOutput(gs)
	}

	return out, nil
}

func (s *GameSystemService) Update(ctx context.Context, in UpdateGameSystemInput) error {
	gs := GameSystem{
		ID:        GameSystemID{in.ID},
		Name:      in.Name,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(ctx, gs)
}

func (s *GameSystemService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, GameSystemID{id})
}

func toGameSystemOutput(gs GameSystem) GameSystemOutput {
	return GameSystemOutput{
		ID:        gs.ID.UUID,
		Name:      gs.Name,
		CreatedAt: gs.CreatedAt,
		UpdatedAt: gs.UpdatedAt,
	}
}
