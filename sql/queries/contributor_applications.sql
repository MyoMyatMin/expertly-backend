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
    ca.contri_app_id,
    ca.user_id,
    ca.expertise_proofs,
    ca.identity_proof,
    ca.initial_submission,
    ca.status,
    ca.created_at,
    ca.reviewed_at,
    u.name AS name,
    u.username AS username
FROM contributor_applications ca
JOIN users u ON ca.user_id = u.user_id
WHERE ca.contri_app_id = $1;

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
    ca.contri_app_id,
    ca.user_id,
    ca.expertise_proofs,
    ca.identity_proof,
    ca.initial_submission,
    ca.status,
    ca.created_at,
    ca.reviewed_at,
    ca.reviewed_by,
    u.name AS name
FROM contributor_applications ca
JOIN users u ON ca.user_id = u.user_id
ORDER BY ca.created_at DESC;
