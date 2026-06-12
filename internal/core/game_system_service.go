package core

import (
	"time"

	"github.com/google/uuid"
)

type gameSystemRepository interface {
	Create(gs GameSystem) (GameSystem, error)
	Get(id GameSystemID) (GameSystem, error)
	GetAll() ([]GameSystem, error)
	Update(gs GameSystem) error
	Delete(id GameSystemID) error
}

type GameSystemService struct {
	repo gameSystemRepository
}

func NewGameSystemService(repo gameSystemRepository) *GameSystemService {
	return &GameSystemService{repo: repo}
}

func (s *GameSystemService) Create(in CreateGameSystemInput) (GameSystemOutput, error) {
	gs := GameSystem{Name: in.Name}
	created, err := s.repo.Create(gs)
	if err != nil {
		return GameSystemOutput{}, err
	}

	return toGameSystemOutput(created), nil
}

func (s *GameSystemService) Get(id string) (GameSystemOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return GameSystemOutput{}, err
	}

	gs, err := s.repo.Get(GameSystemID{uid})
	if err != nil {
		return GameSystemOutput{}, err
	}

	return toGameSystemOutput(gs), nil
}

func (s *GameSystemService) GetAll() ([]GameSystemOutput, error) {
	items, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	out := make([]GameSystemOutput, len(items))
	for i, gs := range items {
		out[i] = toGameSystemOutput(gs)
	}

	return out, nil
}

func (s *GameSystemService) Update(in UpdateGameSystemInput) error {
	uid, err := uuid.Parse(in.ID)
	if err != nil {
		return err
	}

	gs := GameSystem{
		ID:        GameSystemID{uid},
		Name:      in.Name,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(gs)
}

func (s *GameSystemService) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(GameSystemID{uid})
}

func toGameSystemOutput(gs GameSystem) GameSystemOutput {
	return GameSystemOutput{
		ID:        gs.ID.String(),
		Name:      gs.Name,
		CreatedAt: gs.CreatedAt,
		UpdatedAt: gs.UpdatedAt,
	}
}
