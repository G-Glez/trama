-- name: CreateTournament :exec
INSERT INTO tournaments (id, name, player_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?);

-- name: GetTournament :one
SELECT * FROM tournaments WHERE id = ?;

-- name: ListTournaments :many
SELECT * FROM tournaments ORDER BY name;

-- name: UpdateTournament :exec
UPDATE tournaments SET name = ?, player_id = ?, updated_at = ? WHERE id = ?;

-- name: DeleteTournament :exec
DELETE FROM tournaments WHERE id = ?;
