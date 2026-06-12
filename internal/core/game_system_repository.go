package core

import (
	"context"
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
	q gameSystemQuerier
}

func NewGameSystemRepository(q gameSystemQuerier) *GameSystemSQLRepository {
	return &GameSystemSQLRepository{q: q}
}

func (r *GameSystemSQLRepository) Create(gs GameSystem) (GameSystem, error) {
	gs.ID = GameSystemID{uuid.New()}
	now := time.Now()
	gs.CreatedAt = now
	gs.UpdatedAt = now

	err := r.q.CreateGameSystem(context.Background(), coregen.CreateGameSystemParams{
		ID:        gs.ID.String(),
		Name:      gs.Name,
		CreatedAt: gs.CreatedAt,
		UpdatedAt: gs.UpdatedAt,
	})

	return gs, err
}

func (r *GameSystemSQLRepository) Get(id GameSystemID) (GameSystem, error) {
	g, err := r.q.GetGameSystem(context.Background(), id.String())
	if err != nil {
		return GameSystem{}, err
	}

	uid, err := uuid.Parse(g.ID)
	if err != nil {
		return GameSystem{}, err
	}

	return GameSystem{
		ID:        GameSystemID{uid},
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}, nil
}

func (r *GameSystemSQLRepository) GetAll() ([]GameSystem, error) {
	items, err := r.q.ListGameSystems(context.Background())
	if err != nil {
		return nil, err
	}

	list := make([]GameSystem, len(items))
	for i, g := range items {
		uid, err := uuid.Parse(g.ID)
		if err != nil {
			return nil, err
		}
		list[i] = GameSystem{
			ID:        GameSystemID{uid},
			Name:      g.Name,
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		}
	}
	return list, nil
}

func (r *GameSystemSQLRepository) Update(gs GameSystem) error {
	return r.q.UpdateGameSystem(context.Background(), coregen.UpdateGameSystemParams{
		Name:      gs.Name,
		UpdatedAt: gs.UpdatedAt,
		ID:        gs.ID.String(),
	})
}

func (r *GameSystemSQLRepository) Delete(id GameSystemID) error {
	return r.q.DeleteGameSystem(context.Background(), id.String())
}
