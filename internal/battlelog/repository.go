package battlelog

type PlayerRepository interface {
	Create(p *Player) error
	GetByID(id PlayerID) (*Player, error)
	List() ([]Player, error)
	Update(p *Player) error
	Delete(id PlayerID) error
}

type TournamentRepository interface {
	Create(t *Tournament) error
	GetByID(id TournamentID) (*Tournament, error)
	List() ([]Tournament, error)
	Update(t *Tournament) error
	Delete(id TournamentID) error
}

type PairingRepository interface {
	Create(p *Pairing) error
	GetByID(id PairingID) (*Pairing, error)
	ListByRound(roundID string) ([]Pairing, error)
	Update(p *Pairing) error
	Delete(id PairingID) error
}

type BattleResultRepository interface {
	Create(r *BattleResult) error
	GetByID(id BattleResultID) (*BattleResult, error)
	ListByPairing(pairingID PairingID) ([]BattleResult, error)
	Delete(id BattleResultID) error
}
