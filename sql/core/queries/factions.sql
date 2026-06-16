-- name: GetFaction :one
SELECT * FROM factions WHERE id = ?;

-- name: ListFactionsByEdition :many
SELECT * FROM factions WHERE edition_id = ? ORDER BY name;

-- name: CreateFaction :exec
INSERT INTO factions (id, edition_id, game_system_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateFaction :execresult
UPDATE factions SET edition_id = ?, game_system_id = ?, name = ?, updated_at = ? WHERE id = ?;

-- name: DeleteFaction :exec
DELETE FROM factions WHERE id = ?;
