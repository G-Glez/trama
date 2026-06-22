-- name: AddTournamentPlayer :exec
INSERT INTO tournament_players (tournament_id, player_id, team_id) VALUES (?, ?, ?);

-- name: RemoveTournamentPlayer :exec
DELETE FROM tournament_players WHERE tournament_id = ? AND player_id = ?;

-- name: ListTournamentPlayers :many
SELECT * FROM tournament_players WHERE tournament_id = ?;

-- name: DeleteTournamentPlayersByTournament :exec
DELETE FROM tournament_players WHERE tournament_id = ?;
