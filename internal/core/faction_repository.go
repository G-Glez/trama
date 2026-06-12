package core

import (
	"context"
	"time"

	"github.com/google/uuid"

	coregen "trama/internal/gen/core"
)

type factionQuerier interface {
	CreateFaction(ctx context.Context, arg coregen.CreateFactionParams) error
	GetFaction(ctx context.Context, id string) (coregen.Faction, error)
	ListFactionsByEdition(ctx context.Context, editionID string) ([]coregen.Faction, error)
	UpdateFaction(ctx context.Context, arg coregen.UpdateFactionParams) error
	DeleteFaction(ctx context.Context, id string) error
}

type FactionSQLRepository struct {
	q factionQuerier
}

func NewFactionRepository(q factionQuerier) *FactionSQLRepository {
	return &FactionSQLRepository{q: q}
}

func (r *FactionSQLRepository) Create(f Faction) (Faction, error) {
	f.ID = FactionID{uuid.New()}
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now

	err := r.q.CreateFaction(context.Background(), coregen.CreateFactionParams{
		ID:        f.ID.String(),
		EditionID: f.EditionID.String(),
		Name:      f.Name,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	})

	return f, err
}

func (r *FactionSQLRepository) Get(id FactionID) (Faction, error) {
	f, err := r.q.GetFaction(context.Background(), id.String())
	if err != nil {
		return Faction{}, err
	}

	return *mapFaction(&f), nil
}

func (r *FactionSQLRepository) GetAllByEdition(edID EditionID) ([]Faction, error) {
	items, err := r.q.ListFactionsByEdition(context.Background(), edID.String())
	if err != nil {
		return nil, err
	}

	list := make([]Faction, len(items))
	for i, f := range items {
		list[i] = *mapFaction(&f)
	}
	return list, nil
}

func (r *FactionSQLRepository) Update(f Faction) error {
	return r.q.UpdateFaction(context.Background(), coregen.UpdateFactionParams{
		EditionID: f.EditionID.String(),
		Name:      f.Name,
		UpdatedAt: f.UpdatedAt,
		ID:        f.ID.String(),
	})
}

func (r *FactionSQLRepository) Delete(id FactionID) error {
	return r.q.DeleteFaction(context.Background(), id.String())
}

func mapFaction(f *coregen.Faction) *Faction {
	uid, _ := uuid.Parse(f.ID)
	edUID, _ := uuid.Parse(f.EditionID)

	return &Faction{
		ID:        FactionID{uid},
		EditionID: EditionID{edUID},
		Name:      f.Name,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
