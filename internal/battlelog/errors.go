package battlelog

import "errors"

var (
	ErrInvalidPoints      = errors.New("points must be between 0 and 100")
	ErrInvalidWinner      = errors.New("winner must be player1 or player2")
	ErrMatchNotFound      = errors.New("match not found")
	ErrTeamNotFound       = errors.New("team not found")
	ErrTournamentNotFound = errors.New("tournament not found")
)
