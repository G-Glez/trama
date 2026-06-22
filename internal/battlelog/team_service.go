package battlelog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type teamRepository interface {
	Create(ctx context.Context, t Team) (Team, error)
	Get(ctx context.Context, id TeamID) (Team, error)
	List(ctx context.Context) ([]Team, error)
	Update(ctx context.Context, t Team) error
	Delete(ctx context.Context, id TeamID) error
}

type TeamService struct {
	repo teamRepository
}

func NewTeamService(repo teamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) Create(ctx context.Context, in CreateTeamInput) (TeamOutput, error) {
	t := Team{Name: in.Name}

	created, err := s.repo.Create(ctx, t)
	if err != nil {
		return TeamOutput{}, err
	}

	return toTeamOutput(created), nil
}

func (s *TeamService) Get(ctx context.Context, id uuid.UUID) (TeamOutput, error) {
	t, err := s.repo.Get(ctx, TeamID{id})
	if err != nil {
		if errors.Is(err, ErrTeamNotFound) {
			return TeamOutput{}, fmt.Errorf("team %s: %w", id, ErrTeamNotFound)
		}
		return TeamOutput{}, err
	}

	return toTeamOutput(t), nil
}

func (s *TeamService) List(ctx context.Context) ([]TeamOutput, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]TeamOutput, len(items))
	for i, t := range items {
		out[i] = toTeamOutput(t)
	}
	return out, nil
}

func (s *TeamService) Update(ctx context.Context, in UpdateTeamInput) error {
	t := Team{
		ID:        TeamID{in.ID},
		Name:      in.Name,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(ctx, t)
}

func (s *TeamService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, TeamID{id})
}

func toTeamOutput(t Team) TeamOutput {
	return TeamOutput{
		ID:        t.ID.UUID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
