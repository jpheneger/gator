-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
  inserted_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds
  ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users
  ON users.id = inserted_feed_follow.user_id
;

-- name: DeleteFeedFollowss :exec
DELETE
FROM feed_follows
WHERE true
AND user_id IN (SELECT id FROM users WHERE users.name = $1)
AND feed_id IN (SELECT id FROM feeds WHERE feeds.url = $2)
;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, ff.created_at, ff.updated_at, u.name as user_name, f.name as feed_name
FROM feed_follows as ff
INNER JOIN users as u
  ON u.id = ff.user_id
INNER JOIN feeds as f
  ON f.id = ff.feed_id
WHERE 
  ff.user_id IN (SELECT id FROM users WHERE users.name = $1)
;
