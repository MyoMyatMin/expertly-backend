-- name: GetFollowingList :many
SELECT 
    following.following_id,
    users.name,
    users.username
FROM following
JOIN users ON following.following_id = users.user_id
WHERE following.follower_id = $1;

-- name: GetFollowerList :many
SELECT 
    following.follower_id,
    users.name,
    users.username
FROM following
JOIN users ON following.follower_id = users.user_id
WHERE following.following_id = $1;

-- name: CreateFollow :exec
INSERT INTO following (follower_id, following_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteFollow :exec
DELETE FROM following
WHERE follower_id = $1 AND following_id = $2;

-- name: GetFeed :many
SELECT posts.*, users.name, users.username
FROM posts
JOIN following ON posts.user_id = following.following_id
JOIN users ON posts.user_id = users.user_id
WHERE following.follower_id = $1
ORDER BY posts.created_at DESC;

-- name: GetFollwersCount :one
SELECT COUNT(follower_id)
FROM following
WHERE following_id = $1;

-- name: GetFollowingCount :one
SELECT COUNT(following_id)
FROM following
WHERE follower_id = $1;

