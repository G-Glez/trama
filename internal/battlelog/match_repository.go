package battlelog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	battleloggen "trama/internal/gen/battlelog"
	"trama/internal/user"
)

type matchQuerier interface {
	CreateMatch(ctx context.Context, arg battleloggen.CreateMatchParams) error
	GetMatch(ctx context.Context, id string) (battleloggen.Wh40k11thMatch, error)
	ListMatches(ctx context.Context) ([]battleloggen.Wh40k11thMatch, error)
	UpdateMatch(ctx context.Context, arg battleloggen.UpdateMatchParams) error
	DeleteMatch(ctx context.Context, id string) error
}

type MatchSQLRepository struct {
	querier matchQuerier
}

func NewMatchRepository(querier matchQuerier) *MatchSQLRepository {
	return &MatchSQLRepository{querier: querier}
}

func (r *MatchSQLRepository) Create(ctx context.Context, m Match) (Match, error) {
	m.ID = MatchID{uuid.New()}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	err := r.querier.CreateMatch(ctx, battleloggen.CreateMatchParams{
		ID:        m.ID.String(),
		Player1ID: m.Player1ID.String(),
		Player2ID: m.Player2ID.String(),
		Points1:   int64(m.Points1),
		Points2:   int64(m.Points2),
		WinnerID:  m.WinnerID.String(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	})
	if err != nil {
		return Match{}, fmt.Errorf("database error: %w", err)
	}

	return m, nil
}

func (r *MatchSQLRepository) Get(ctx context.Context, id MatchID) (Match, error) {
	m, err := r.querier.GetMatch(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Match{}, ErrMatchNotFound
		}
		return Match{}, fmt.Errorf("database error: %w", err)
	}

	return toDomainMatch(m), nil
}

func (r *MatchSQLRepository) List(ctx context.Context) ([]Match, error) {
	items, err := r.querier.ListMatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	list := make([]Match, len(items))
	for i, m := range items {
		list[i] = toDomainMatch(m)
	}
	return list, nil
}

func (r *MatchSQLRepository) Update(ctx context.Context, m Match) error {
	err := r.querier.UpdateMatch(ctx, battleloggen.UpdateMatchParams{
		Player1ID: m.Player1ID.String(),
		Player2ID: m.Player2ID.String(),
		Points1:   int64(m.Points1),
		Points2:   int64(m.Points2),
		WinnerID:  m.WinnerID.String(),
		UpdatedAt: m.UpdatedAt,
		ID:        m.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *MatchSQLRepository) Delete(ctx context.Context, id MatchID) error {
	err := r.querier.DeleteMatch(ctx, id.String())
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func toDomainMatch(m battleloggen.Wh40k11thMatch) Match {
	return Match{
		ID:        MatchID{uuid.MustParse(m.ID)},
		Player1ID: user.NewUserID(uuid.MustParse(m.Player1ID)),
		Player2ID: user.NewUserID(uuid.MustParse(m.Player2ID)),
		Points1:   int(m.Points1),
		Points2:   int(m.Points2),
		WinnerID:  user.NewUserID(uuid.MustParse(m.WinnerID)),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
