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

type gameSystemQuerier interface {
	CreateGameSystem(ctx context.Context, arg coregen.CreateGameSystemParams) error
	GetGameSystem(ctx context.Context, id string) (coregen.GameSystem, error)
	ListGameSystems(ctx context.Context) ([]coregen.GameSystem, error)
	UpdateGameSystem(ctx context.Context, arg coregen.UpdateGameSystemParams) error
	DeleteGameSystem(ctx context.Context, id string) error
}

type GameSystemSQLRepository struct {
	querier gameSystemQuerier
}

func NewGameSystemRepository(querier gameSystemQuerier) *GameSystemSQLRepository {
	return &GameSystemSQLRepository{querier: querier}
}

func (r *GameSystemSQLRepository) Create(ctx context.Context, gs GameSystem) (GameSystem, error) {
	gs.ID = GameSystemID{uuid.New()}
	now := time.Now()
	gs.CreatedAt = now
	gs.UpdatedAt = now

	err := r.querier.CreateGameSystem(ctx, coregen.CreateGameSystemParams{
		ID:        gs.ID.String(),
		Name:      gs.Name,
		CreatedAt: gs.CreatedAt,
		UpdatedAt: gs.UpdatedAt,
	})
	if err != nil {
		return GameSystem{}, fmt.Errorf("%w: %w", ErrDB, err)
	}

	return gs, nil
}

func (r *GameSystemSQLRepository) Get(ctx context.Context, id GameSystemID) (GameSystem, error) {
	g, err := r.querier.GetGameSystem(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GameSystem{}, ErrNotFound
		}
		return GameSystem{}, fmt.Errorf("%w: %w", ErrDB, err)
	}

	return toDomainGameSystem(g)
}

func (r *GameSystemSQLRepository) GetAll(ctx context.Context) ([]GameSystem, error) {
	items, err := r.querier.ListGameSystems(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDB, err)
	}

	list := make([]GameSystem, len(items))
	for i, g := range items {
		list[i], err = toDomainGameSystem(g)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (r *GameSystemSQLRepository) Update(ctx context.Context, gs GameSystem) error {
	err := r.querier.UpdateGameSystem(ctx, coregen.UpdateGameSystemParams{
		Name:      gs.Name,
		UpdatedAt: gs.UpdatedAt,
		ID:        gs.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDB, err)
	}
	return nil
}

func (r *GameSystemSQLRepository) Delete(ctx context.Context, id GameSystemID) error {
	err := r.querier.DeleteGameSystem(ctx, id.String())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDB, err)
	}
	return nil
}

func toDomainGameSystem(g coregen.GameSystem) (GameSystem, error) {
	uid, err := uuid.Parse(g.ID)
	if err != nil {
		return GameSystem{}, fmt.Errorf("%w: invalid uuid: %w", ErrDataCorruption, err)
	}

	return GameSystem{
		ID:        GameSystemID{uid},
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}, nil
}
