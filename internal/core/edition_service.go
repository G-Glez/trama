package core

import (
	"time"

	"github.com/google/uuid"
)

type editionRepository interface {
	Create(e Edition) (Edition, error)
	Get(id EditionID) (Edition, error)
	GetAllByGameSystem(gsID GameSystemID) ([]Edition, error)
	Update(e Edition) error
	Delete(id EditionID) error
}

type EditionService struct {
	repo editionRepository
}

func NewEditionService(repo editionRepository) *EditionService {
	return &EditionService{repo: repo}
}

func (s *EditionService) Create(in CreateEditionInput) (EditionOutput, error) {
	gsID, err := uuid.Parse(in.GameSystemID)
	if err != nil {
		return EditionOutput{}, err
	}

	e := Edition{
		GameSystemID: GameSystemID{gsID},
		Name:         in.Name,
		Version:      in.Version,
	}

	created, err := s.repo.Create(e)
	if err != nil {
		return EditionOutput{}, err
	}

	return toEditionOutput(created), nil
}

func (s *EditionService) Get(id string) (EditionOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return EditionOutput{}, err
	}

	e, err := s.repo.Get(EditionID{uid})
	if err != nil {
		return EditionOutput{}, err
	}

	return toEditionOutput(e), nil
}

func (s *EditionService) GetAllByGameSystem(gsID string) ([]EditionOutput, error) {
	uid, err := uuid.Parse(gsID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetAllByGameSystem(GameSystemID{uid})
	if err != nil {
		return nil, err
	}

	out := make([]EditionOutput, len(items))
	for i, e := range items {
		out[i] = toEditionOutput(e)
	}

	return out, nil
}

func (s *EditionService) Update(in UpdateEditionInput) error {
	uid, err := uuid.Parse(in.ID)
	if err != nil {
		return err
	}

	gsID, err := uuid.Parse(in.GameSystemID)
	if err != nil {
		return err
	}

	e := Edition{
		ID:           EditionID{uid},
		GameSystemID: GameSystemID{gsID},
		Name:         in.Name,
		Version:      in.Version,
		UpdatedAt:    time.Now(),
	}

	return s.repo.Update(e)
}

func (s *EditionService) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(EditionID{uid})
}

func toEditionOutput(e Edition) EditionOutput {
	return EditionOutput{
		ID:           e.ID.String(),
		GameSystemID: e.GameSystemID.String(),
		Name:         e.Name,
		Version:      e.Version,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
