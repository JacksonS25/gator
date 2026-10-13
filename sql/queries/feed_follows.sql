-- name: CreateFeedFollow :one
INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowForUser :many
SELECT
    id,
    created_at,
    updated_at,
    user_id,
    feed_id
FROM feed_follows
WHERE user_id = $1;

-- name: RemoveFeedFollow :one
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2
RETURNING *;