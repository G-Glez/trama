package core

import (
	"time"

	"github.com/google/uuid"
)

type factionRepository interface {
	Create(f Faction) (Faction, error)
	Get(id FactionID) (Faction, error)
	GetAllByEdition(edID EditionID) ([]Faction, error)
	Update(f Faction) error
	Delete(id FactionID) error
}

type FactionService struct {
	repo factionRepository
}

func NewFactionService(repo factionRepository) *FactionService {
	return &FactionService{repo: repo}
}

func (s *FactionService) Create(in CreateFactionInput) (FactionOutput, error) {
	edID, err := uuid.Parse(in.EditionID)
	if err != nil {
		return FactionOutput{}, err
	}

	f := Faction{
		EditionID: EditionID{edID},
		Name:      in.Name,
	}

	created, err := s.repo.Create(f)
	if err != nil {
		return FactionOutput{}, err
	}

	return toFactionOutput(created), nil
}

func (s *FactionService) Get(id string) (FactionOutput, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return FactionOutput{}, err
	}

	f, err := s.repo.Get(FactionID{uid})
	if err != nil {
		return FactionOutput{}, err
	}

	return toFactionOutput(f), nil
}

func (s *FactionService) GetAllByEdition(edID string) ([]FactionOutput, error) {
	uid, err := uuid.Parse(edID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetAllByEdition(EditionID{uid})
	if err != nil {
		return nil, err
	}

	out := make([]FactionOutput, len(items))
	for i, f := range items {
		out[i] = toFactionOutput(f)
	}

	return out, nil
}

func (s *FactionService) Update(in UpdateFactionInput) error {
	uid, err := uuid.Parse(in.ID)
	if err != nil {
		return err
	}

	edID, err := uuid.Parse(in.EditionID)
	if err != nil {
		return err
	}

	f := Faction{
		ID:        FactionID{uid},
		EditionID: EditionID{edID},
		Name:      in.Name,
		UpdatedAt: time.Now(),
	}

	return s.repo.Update(f)
}

func (s *FactionService) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(FactionID{uid})
}

func toFactionOutput(f Faction) FactionOutput {
	return FactionOutput{
		ID:        f.ID.String(),
		EditionID: f.EditionID.String(),
		Name:      f.Name,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
