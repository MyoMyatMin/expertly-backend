-- name: CreateContributor :one
INSERT INTO Contributors (
    user_id,
    expertise_fields
) VALUES (
    $1,
    $2
)
RETURNING user_id, expertise_fields, created_at;

-- name: GetContributorByUserId :one
SELECT 
    user_id,
    expertise_fields,
    created_at
FROM Contributors
WHERE user_id = $1;

-- name: UpdateContributorExpertiseFields :exec
UPDATE Contributors
SET expertise_fields = $2
WHERE user_id = $1;

-- name: DeleteContributor :exec
DELETE FROM Contributors
WHERE user_id = $1;

-- name: ListAllContributors :many
SELECT 
    user_id,
    expertise_fields,
    created_at
FROM Contributors;

-- name: SearchContributorsByExpertiseField :many
SELECT 
    user_id,
    expertise_fields,
    created_at
FROM Contributors
WHERE $1 = ANY(expertise_fields);