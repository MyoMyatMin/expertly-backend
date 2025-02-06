-- name: ApplyContributorApplication :one
INSERT INTO contributor_applications(
    contri_app_id,
    user_id,
    expertise_proofs,
    identity_proof,
    initial_submission
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING contri_app_id, user_id, expertise_proofs, identity_proof, initial_submission, status, created_at, reviewed_at;

-- name: GetContributorApplication :one
SELECT 
    contri_app_id,
    user_id,
    expertise_proofs,
    identity_proof,
    initial_submission,
    status,
    created_at,
    reviewed_at
FROM contributor_applications
WHERE contri_app_id = $1;

-- name: UpdateContributorApplication :one
UPDATE contributor_applications
SET
    status = $2,
    reviewed_at = CURRENT_TIMESTAMP,
    reviewed_by = $3
WHERE contri_app_id = $1
RETURNING contri_app_id, user_id, expertise_proofs, identity_proof, initial_submission, status, created_at, reviewed_at;

-- name: ListContributorApplications :many
SELECT 
    contri_app_id,
    user_id,
    expertise_proofs,
    identity_proof,
    initial_submission,
    status,
    created_at,
    reviewed_at
FROM contributor_applications
ORDER BY created_at DESC;
