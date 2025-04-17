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

-- name: GetFeedbyUrl :one
SELECT id,created_at,updated_at,name,url,user_id
FROM feeds
WHERE url = $1;

-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id,created_at,updated_at,user_id,feed_id)
    VALUES(
    $1,
    $2,
    $3,
    $4,
    $5    
    )
    RETURNING *
)
SELECT
inserted_feed_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM inserted_feed_follow
INNER JOIN  feeds
ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users
ON inserted_feed_follow.user_id = users.id;