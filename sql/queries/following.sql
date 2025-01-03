-- name: GetFollowingList :many
SELECT followee_id, followed_at
FROM following
WHERE follower_id = $1;

-- name: GetFollowerList :many
SELECT follower_id, followed_at
FROM following
WHERE followee_id = $1;

-- name: CreateFollow :exec
INSERT INTO following (follower_id, followee_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteFollow :exec
DELETE FROM following
WHERE follower_id = $1 AND followee_id = $2;

-- name: GetFeed :many
SELECT posts.*
FROM posts
JOIN following ON posts.user_id = following.followee_id
WHERE following.follower_id = $1
ORDER BY posts.created_at DESC;
