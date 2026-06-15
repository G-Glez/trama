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

type editionQuerier interface {
	CreateEdition(ctx context.Context, arg coregen.CreateEditionParams) error
	GetEdition(ctx context.Context, id string) (coregen.Edition, error)
	ListEditionsByGameSystem(ctx context.Context, gameSystemID string) ([]coregen.Edition, error)
	UpdateEdition(ctx context.Context, arg coregen.UpdateEditionParams) error
	DeleteEdition(ctx context.Context, id string) error
}

type EditionSQLRepository struct {
	querier editionQuerier
}

func NewEditionRepository(querier editionQuerier) *EditionSQLRepository {
	return &EditionSQLRepository{querier: querier}
}

func (r *EditionSQLRepository) Create(ctx context.Context, e Edition) (Edition, error) {
	e.ID = EditionID{uuid.New()}
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now

	err := r.querier.CreateEdition(ctx, coregen.CreateEditionParams{
		ID:           e.ID.String(),
		GameSystemID: e.GameSystemID.String(),
		Name:         e.Name,
		Version:      e.Version,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	})
	if err != nil {
		return Edition{}, fmt.Errorf("%w: %w", ErrDB, err)
	}

	return e, nil
}

func (r *EditionSQLRepository) Get(ctx context.Context, id EditionID) (Edition, error) {
	ed, err := r.querier.GetEdition(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Edition{}, ErrNotFound
		}
		return Edition{}, fmt.Errorf("%w: %w", ErrDB, err)
	}

	return toDomainEdition(ed)
}

func (r *EditionSQLRepository) GetAllByGameSystem(ctx context.Context, gsID GameSystemID) ([]Edition, error) {
	items, err := r.querier.ListEditionsByGameSystem(ctx, gsID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDB, err)
	}

	list := make([]Edition, len(items))
	for i, ed := range items {
		list[i], err = toDomainEdition(ed)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (r *EditionSQLRepository) Update(ctx context.Context, e Edition) error {
	err := r.querier.UpdateEdition(ctx, coregen.UpdateEditionParams{
		GameSystemID: e.GameSystemID.String(),
		Name:         e.Name,
		Version:      e.Version,
		UpdatedAt:    e.UpdatedAt,
		ID:           e.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDB, err)
	}
	return nil
}

func (r *EditionSQLRepository) Delete(ctx context.Context, id EditionID) error {
	err := r.querier.DeleteEdition(ctx, id.String())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDB, err)
	}
	return nil
}

func toDomainEdition(ed coregen.Edition) (Edition, error) {
	uid, err := uuid.Parse(ed.ID)
	if err != nil {
		return Edition{}, fmt.Errorf("%w: invalid uuid: %w", ErrDataCorruption, err)
	}

	gsUID, err := uuid.Parse(ed.GameSystemID)
	if err != nil {
		return Edition{}, fmt.Errorf("%w: invalid uuid: %w", ErrDataCorruption, err)
	}

	return Edition{
		ID:           EditionID{uid},
		GameSystemID: GameSystemID{gsUID},
		Name:         ed.Name,
		Version:      ed.Version,
		CreatedAt:    ed.CreatedAt,
		UpdatedAt:    ed.UpdatedAt,
	}, nil
}
