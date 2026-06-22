-- name: CreateMatch :exec
INSERT INTO wh40k_11th_matches (id, player1_id, player2_id, points1, points2, winner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetMatch :one
SELECT * FROM wh40k_11th_matches WHERE id = ?;

-- name: ListMatches :many
SELECT * FROM wh40k_11th_matches ORDER BY created_at DESC;

-- name: UpdateMatch :exec
UPDATE wh40k_11th_matches SET player1_id = ?, player2_id = ?, points1 = ?, points2 = ?, winner_id = ?, updated_at = ? WHERE id = ?;

-- name: DeleteMatch :exec
DELETE FROM wh40k_11th_matches WHERE id = ?;
