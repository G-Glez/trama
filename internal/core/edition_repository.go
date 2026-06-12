package core

import (
	"context"
	"time"

	"github.com/google/uuid"

	coregen "trama/internal/gen/core"
)

type editionQuerier interface {
	CreateEdition(ctx context.Context, arg coregen.CreateEditionParams) error
	GetEdition(ctx context.Context, id string) (coregen.Edition, error)
	ListEditionsByGameSystem(ctx context.Context, gameSystemID string) ([]coregen.Edition, error)
	UpdateEdition(ctx context.Context, arg coregen.UpdateEditionParams) error
	DeleteEdition(ctx context.Context, id string) error
}

type EditionSQLRepository struct {
	q editionQuerier
}

func NewEditionRepository(q editionQuerier) *EditionSQLRepository {
	return &EditionSQLRepository{q: q}
}

func (r *EditionSQLRepository) Create(e Edition) (Edition, error) {
	e.ID = EditionID{uuid.New()}
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now

	err := r.q.CreateEdition(context.Background(), coregen.CreateEditionParams{
		ID:           e.ID.String(),
		GameSystemID: e.GameSystemID.String(),
		Name:         e.Name,
		Version:      e.Version,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	})

	return e, err
}

func (r *EditionSQLRepository) Get(id EditionID) (Edition, error) {
	ed, err := r.q.GetEdition(context.Background(), id.String())
	if err != nil {
		return Edition{}, err
	}

	return *mapEdition(&ed), nil
}

func (r *EditionSQLRepository) GetAllByGameSystem(gsID GameSystemID) ([]Edition, error) {
	items, err := r.q.ListEditionsByGameSystem(context.Background(), gsID.String())
	if err != nil {
		return nil, err
	}

	list := make([]Edition, len(items))
	for i, ed := range items {
		list[i] = *mapEdition(&ed)
	}
	return list, nil
}

func (r *EditionSQLRepository) Update(e Edition) error {
	return r.q.UpdateEdition(context.Background(), coregen.UpdateEditionParams{
		GameSystemID: e.GameSystemID.String(),
		Name:         e.Name,
		Version:      e.Version,
		UpdatedAt:    e.UpdatedAt,
		ID:           e.ID.String(),
	})
}

func (r *EditionSQLRepository) Delete(id EditionID) error {
	return r.q.DeleteEdition(context.Background(), id.String())
}

func mapEdition(ed *coregen.Edition) *Edition {
	uid, _ := uuid.Parse(ed.ID)
	gsUID, _ := uuid.Parse(ed.GameSystemID)

	return &Edition{
		ID:           EditionID{uid},
		GameSystemID: GameSystemID{gsUID},
		Name:         ed.Name,
		Version:      ed.Version,
		CreatedAt:    ed.CreatedAt,
		UpdatedAt:    ed.UpdatedAt,
	}
}
