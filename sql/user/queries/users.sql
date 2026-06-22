-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash, default_faction, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY username;

-- name: UpdateUser :exec
UPDATE users SET username = ?, email = ?, password_hash = ?, default_faction = ?, updated_at = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
