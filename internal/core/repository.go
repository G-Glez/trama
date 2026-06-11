package core

type GameSystemRepository interface {
	Create(gs *GameSystem) error
	GetByID(id GameSystemID) (*GameSystem, error)
	List() ([]GameSystem, error)
	Update(gs *GameSystem) error
	Delete(id GameSystemID) error
}

type EditionRepository interface {
	Create(e *Edition) error
	GetByID(id EditionID) (*Edition, error)
	ListByGameSystem(gsID GameSystemID) ([]Edition, error)
	Update(e *Edition) error
	Delete(id EditionID) error
}

type FactionRepository interface {
	Create(f *Faction) error
	GetByID(id FactionID) (*Faction, error)
	ListByEdition(edID EditionID) ([]Faction, error)
	Update(f *Faction) error
	Delete(id FactionID) error
}
