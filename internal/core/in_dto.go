package core

type CreateGameSystemInput struct {
	Name string
}

type UpdateGameSystemInput struct {
	ID   string
	Name string
}

type CreateEditionInput struct {
	GameSystemID string
	Name         string
	Version      string
}

type UpdateEditionInput struct {
	ID           string
	GameSystemID string
	Name         string
	Version      string
}

type CreateFactionInput struct {
	EditionID string
	Name      string
}

type UpdateFactionInput struct {
	ID         string
	EditionID  string
	Name       string
}
