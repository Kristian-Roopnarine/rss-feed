-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT *
FROM feeds;

-- name: GetNextFeedsToFetch :many
select *
from feeds
order by last_fetched_at asc
limit @number_of_results;

-- name: MarkFeedFetched :one
UPDATE feeds
SET updated_at = @new_time, last_fetched_at = @new_time
WHERE id = @feed_id
RETURNING *;