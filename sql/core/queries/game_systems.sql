-- name: GetGameSystem :one
SELECT * FROM game_systems WHERE id = ?;

-- name: ListGameSystems :many
SELECT * FROM game_systems ORDER BY name;

-- name: CreateGameSystem :exec
INSERT INTO game_systems (id, name, created_at, updated_at) VALUES (?, ?, ?, ?);

-- name: UpdateGameSystem :exec
UPDATE game_systems SET name = ?, updated_at = ? WHERE id = ?;

-- name: DeleteGameSystem :exec
DELETE FROM game_systems WHERE id = ?;
