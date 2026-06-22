-- name: CreateTeam :exec
INSERT INTO teams (id, name, created_at, updated_at) VALUES (?, ?, ?, ?);

-- name: GetTeam :one
SELECT * FROM teams WHERE id = ?;

-- name: ListTeams :many
SELECT * FROM teams ORDER BY name;

-- name: UpdateTeam :exec
UPDATE teams SET name = ?, updated_at = ? WHERE id = ?;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = ?;
