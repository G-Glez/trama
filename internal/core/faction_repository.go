package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	coregen "trama/internal/gen/core"
)

type factionQuerier interface {
	CreateFaction(ctx context.Context, arg coregen.CreateFactionParams) error
	GetFaction(ctx context.Context, id string) (coregen.Faction, error)
	ListFactionsByEdition(ctx context.Context, editionID string) ([]coregen.Faction, error)
	UpdateFaction(ctx context.Context, arg coregen.UpdateFactionParams) (sql.Result, error)
	DeleteFaction(ctx context.Context, id string) error
}

type FactionSQLRepository struct {
	querier factionQuerier
}

func NewFactionRepository(querier factionQuerier) *FactionSQLRepository {
	return &FactionSQLRepository{querier: querier}
}

func (r *FactionSQLRepository) Create(ctx context.Context, f Faction) (Faction, error) {
	f.ID = FactionID{uuid.New()}
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now

	err := r.querier.CreateFaction(ctx, coregen.CreateFactionParams{
		ID:           f.ID.String(),
		EditionID:    f.EditionID.String(),
		GameSystemID: f.GameSystemID.String(),
		Name:         f.Name,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	})
	if err != nil {
		return Faction{}, fmt.Errorf("%w: %w", ErrDB, err)
	}

	return f, nil
}

func (r *FactionSQLRepository) Get(ctx context.Context, id FactionID) (Faction, error) {
	f, err := r.querier.GetFaction(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Faction{}, ErrNotFound
		}
		return Faction{}, fmt.Errorf("%w: %w", ErrDB, err)
	}

	return toDomainFaction(f)
}

func (r *FactionSQLRepository) GetAllByEdition(ctx context.Context, edID EditionID) ([]Faction, error) {
	items, err := r.querier.ListFactionsByEdition(ctx, edID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDB, err)
	}

	list := make([]Faction, len(items))
	for i, f := range items {
		list[i], err = toDomainFaction(f)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (r *FactionSQLRepository) Update(ctx context.Context, f Faction) error {
	res, err := r.querier.UpdateFaction(ctx, coregen.UpdateFactionParams{
		EditionID:    f.EditionID.String(),
		GameSystemID: f.GameSystemID.String(),
		Name:         f.Name,
		UpdatedAt:    f.UpdatedAt,
		ID:           f.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDB, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *FactionSQLRepository) Delete(ctx context.Context, id FactionID) error {
	err := r.querier.DeleteFaction(ctx, id.String())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDB, err)
	}
	return nil
}

func toDomainFaction(f coregen.Faction) (Faction, error) {
	uid, err := uuid.Parse(f.ID)
	if err != nil {
		return Faction{}, fmt.Errorf("%w: invalid uuid: %w", ErrDataCorruption, err)
	}

	edUID, err := uuid.Parse(f.EditionID)
	if err != nil {
		return Faction{}, fmt.Errorf("%w: invalid uuid: %w", ErrDataCorruption, err)
	}

	gsUID, err := uuid.Parse(f.GameSystemID)
	if err != nil {
		return Faction{}, fmt.Errorf("%w: invalid uuid: %w", ErrDataCorruption, err)
	}

	return Faction{
		ID:           FactionID{uid},
		EditionID:    EditionID{edUID},
		GameSystemID: GameSystemID{gsUID},
		Name:         f.Name,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}, nil
}
