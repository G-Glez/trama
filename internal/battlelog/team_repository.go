package battlelog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	battleloggen "trama/internal/gen/battlelog"
)

type teamQuerier interface {
	CreateTeam(ctx context.Context, arg battleloggen.CreateTeamParams) error
	GetTeam(ctx context.Context, id string) (battleloggen.Team, error)
	ListTeams(ctx context.Context) ([]battleloggen.Team, error)
	UpdateTeam(ctx context.Context, arg battleloggen.UpdateTeamParams) error
	DeleteTeam(ctx context.Context, id string) error
}

type TeamSQLRepository struct {
	querier teamQuerier
}

func NewTeamRepository(querier teamQuerier) *TeamSQLRepository {
	return &TeamSQLRepository{querier: querier}
}

func (r *TeamSQLRepository) Create(ctx context.Context, t Team) (Team, error) {
	t.ID = TeamID{uuid.New()}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	err := r.querier.CreateTeam(ctx, battleloggen.CreateTeamParams{
		ID:        t.ID.String(),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	})
	if err != nil {
		return Team{}, fmt.Errorf("database error: %w", err)
	}

	return t, nil
}

func (r *TeamSQLRepository) Get(ctx context.Context, id TeamID) (Team, error) {
	t, err := r.querier.GetTeam(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Team{}, ErrTeamNotFound
		}
		return Team{}, fmt.Errorf("database error: %w", err)
	}

	return toDomainTeam(t), nil
}

func (r *TeamSQLRepository) List(ctx context.Context) ([]Team, error) {
	items, err := r.querier.ListTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	list := make([]Team, len(items))
	for i, t := range items {
		list[i] = toDomainTeam(t)
	}
	return list, nil
}

func (r *TeamSQLRepository) Update(ctx context.Context, t Team) error {
	err := r.querier.UpdateTeam(ctx, battleloggen.UpdateTeamParams{
		Name:      t.Name,
		UpdatedAt: t.UpdatedAt,
		ID:        t.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *TeamSQLRepository) Delete(ctx context.Context, id TeamID) error {
	err := r.querier.DeleteTeam(ctx, id.String())
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func toDomainTeam(t battleloggen.Team) Team {
	return Team{
		ID:        TeamID{uuid.MustParse(t.ID)},
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
