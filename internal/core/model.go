package core

import (
	"time"

	"github.com/google/uuid"
)

var GameSystems = GameSystemCatalog{
	Warhammer40k: Warhammer40kSystem{
		ID:   GameSystemID{uuid.MustParse("be9d0e8f-c31d-4609-a4c1-5d3496e2bd35")},
		Name: "Warhammer 40,000",
		Editions: Warhammer40kEditions{
			Edition11: Warhammer40kEdition11{
				ID:   EditionID{uuid.MustParse("9577d52c-8ffa-4e9b-9a15-1b46a507d25b")},
				Name: "11th Edition",
				Patches: Warhammer40kEdition11Patches{
					Launch: Patch{
						ID:   PatchID{uuid.MustParse("11111111-1111-1111-1111-111111111111")},
						Name: "Launch",
					},
				},
			},
		},
	},
}

type GameSystemCatalog struct {
	Warhammer40k Warhammer40kSystem
}

type Warhammer40kSystem struct {
	ID       GameSystemID
	Name     string
	Editions Warhammer40kEditions
}

type Warhammer40kEditions struct {
	Edition11 Warhammer40kEdition11
}

type Warhammer40kEdition11 struct {
	ID      EditionID
	Name    string
	Patches Warhammer40kEdition11Patches
}

type Warhammer40kEdition11Patches struct {
	Launch Patch
}

type GameSystemID struct{ uuid.UUID }

type EditionID struct{ uuid.UUID }

type PatchID struct{ uuid.UUID }

type Patch struct {
	ID   PatchID
	Name string
}

type FactionID struct{ uuid.UUID }

type Faction struct {
	ID           FactionID
	EditionID    EditionID
	GameSystemID GameSystemID
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
