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
TRUNCATE TABLE users, feeds,feed_follows, posts;

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

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name, users.name
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = (
    SELECT id FROM users WHERE users.name = $1
)
AND feed_id = (
    SELECT id FROM feeds WHERE url = $2
);

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = (NOW()), last_fetched_at = (NOW())
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT name, url, id 
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
FETCH FIRST 1 ROW ONLY;

-- name: CreatePost :one
INSERT INTO posts (id,created_at,updated_at,title,url,description,published_at,feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.title, p.url, p.published_at, p.description
FROM posts p
JOIN feeds f ON p.feed_id = f.id
JOIN feed_follows fw ON f.id = fw.feed_id
WHERE fw.user_id = $1
ORDER BY published_at DESC
LIMIT $2;

