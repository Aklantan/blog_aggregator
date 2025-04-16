-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT id, created_at, updated_at, name FROM users
WHERE name = $1;


-- name: ResetUsers :exec
TRUNCATE TABLE users, feeds;

-- name: GetUsers :many
SELECT id, created_at, updated_at, name FROM users;

-- name: CreateFeed :one
INSERT INTO feeds (id,created_at,updated_at,name,url,user_id)
VALUES (
     $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT user_id, name, url FROM feeds;

-- name: GetFeedUser :one
SELECT name FROM users
WHERE id = $1;