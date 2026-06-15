package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type editionRepository interface {
	Create(ctx context.Context, e Edition) (Edition, error)
	Get(ctx context.Context, id EditionID) (Edition, error)
	GetAllByGameSystem(ctx context.Context, gsID GameSystemID) ([]Edition, error)
	Update(ctx context.Context, e Edition) error
	Delete(ctx context.Context, id EditionID) error
}

type EditionService struct {
	repo editionRepository
}

func NewEditionService(repo editionRepository) *EditionService {
	return &EditionService{repo: repo}
}

func (s *EditionService) Create(ctx context.Context, in CreateEditionInput) (EditionOutput, error) {
	e := Edition{
		GameSystemID: GameSystemID{in.GameSystemID},
		Name:         in.Name,
		Version:      in.Version,
	}

	created, err := s.repo.Create(ctx, e)
	if err != nil {
		return EditionOutput{}, err
	}

	return toEditionOutput(created), nil
}

func (s *EditionService) Get(ctx context.Context, id uuid.UUID) (EditionOutput, error) {
	e, err := s.repo.Get(ctx, EditionID{id})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return EditionOutput{}, fmt.Errorf("edition %s: %w", id, ErrNotFound)
		}
		return EditionOutput{}, err
	}

	return toEditionOutput(e), nil
}

func (s *EditionService) GetAllByGameSystem(ctx context.Context, gsID uuid.UUID) ([]EditionOutput, error) {
	items, err := s.repo.GetAllByGameSystem(ctx, GameSystemID{gsID})
	if err != nil {
		return nil, err
	}

	out := make([]EditionOutput, len(items))
	for i, e := range items {
		out[i] = toEditionOutput(e)
	}

	return out, nil
}

func (s *EditionService) Update(ctx context.Context, in UpdateEditionInput) error {
	e := Edition{
		ID:           EditionID{in.ID},
		GameSystemID: GameSystemID{in.GameSystemID},
		Name:         in.Name,
		Version:      in.Version,
		UpdatedAt:    time.Now(),
	}

	return s.repo.Update(ctx, e)
}

func (s *EditionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, EditionID{id})
}

func toEditionOutput(e Edition) EditionOutput {
	return EditionOutput{
		ID:           e.ID.UUID,
		GameSystemID: e.GameSystemID.UUID,
		Name:         e.Name,
		Version:      e.Version,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
