-- name: GetEdition :one
SELECT * FROM editions WHERE id = ?;

-- name: ListEditionsByGameSystem :many
SELECT * FROM editions WHERE game_system_id = ? ORDER BY name;

-- name: CreateEdition :exec
INSERT INTO editions (id, game_system_id, name, version, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateEdition :exec
UPDATE editions SET game_system_id = ?, name = ?, version = ?, updated_at = ? WHERE id = ?;

-- name: DeleteEdition :exec
DELETE FROM editions WHERE id = ?;
