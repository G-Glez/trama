package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type factionRepository interface {
	Create(ctx context.Context, f Faction) (Faction, error)
	Get(ctx context.Context, id FactionID) (Faction, error)
	GetAllByEdition(ctx context.Context, edID EditionID) ([]Faction, error)
	Update(ctx context.Context, f Faction) error
	Delete(ctx context.Context, id FactionID) error
}

type FactionService struct {
	repo factionRepository
}

func NewFactionService(repo factionRepository) *FactionService {
	return &FactionService{repo: repo}
}

func (s *FactionService) Create(ctx context.Context, in CreateFactionInput) (FactionOutput, error) {
	f := Faction{
		EditionID: EditionID{in.EditionID},
		Name:      in.Name,
	}

	created, err := s.repo.Create(ctx, f)
	if err != nil {
		return FactionOutput{}, err
	}

	return toFactionOutput(created), nil
}

func (s *FactionService) Get(ctx context.Context, id uuid.UUID) (FactionOutput, error) {
	f, err := s.repo.Get(ctx, FactionID{id})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return FactionOutput{}, fmt.Errorf("faction %s: %w", id, ErrNotFound)
		}
		return FactionOutput{}, err
	}

	return toFactionOutput(f), nil
}

func (s *FactionService) GetAllByEdition(ctx context.Context, edID uuid.UUID) ([]FactionOutput, error) {
	items, err := s.repo.GetAllByEdition(ctx, EditionID{edID})
	if err != nil {
		return nil, err
	}

	out := make([]FactionOutput, len(items))
	for i, f := range items {
		out[i] = toFactionOutput(f)
	}

	return out, nil
}

func (s *FactionService) Update(ctx context.Context, in UpdateFactionInput) error {
	f := Faction{
		ID:        FactionID{in.ID},
		EditionID: EditionID{in.EditionID},
		Name:      in.Name,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(ctx, f)
}

func (s *FactionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, FactionID{id})
}

func toFactionOutput(f Faction) FactionOutput {
	return FactionOutput{
		ID:        f.ID.UUID,
		EditionID: f.EditionID.UUID,
		Name:      f.Name,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
