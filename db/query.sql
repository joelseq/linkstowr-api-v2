-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES (?, ?)
RETURNING id, username;

-- name: GetUser :one
SELECT id, username, password FROM users
WHERE username = ?;

-- name: CreateToken :exec
INSERT INTO tokens (token_hash, name, short_token, user_id)
VALUES (?, ?, ?, ?);

-- name: GetToken :one
SELECT id, name, user_id FROM tokens
WHERE token_hash = ?;

-- name: ListTokens :many
SELECT id, name, short_token FROM tokens
WHERE user_id = ?;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE id = ? AND user_id = ?;

-- name: CreateLink :one
INSERT INTO links (url, title, note, user_id, tags)
VALUES (?, ?, ?, ?, ?)
RETURNING id, url;

-- name: ListLinks :many
SELECT url, title, note, bookmarked_at, tags FROM links
WHERE user_id = ?;

-- name: ClearLinks :exec
DELETE FROM links
WHERE user_id = ?;
