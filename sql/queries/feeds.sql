-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds 
ORDER BY 
  (
    CASE WHEN last_fetched IS NULL THEN 1 ELSE 0 END
  ) DESC, 
  last_fetched ASC
LIMIT $1;
