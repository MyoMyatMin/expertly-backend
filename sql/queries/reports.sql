-- name: CreateReport :one
INSERT INTO reports (report_id, reported_by, target_post_id, target_user_id, target_comment_id, reason) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateReportStatus :one
UPDATE reports SET status = $1, reviewedby = $2, reviewed_at = CURRENT_TIMESTAMP, suspend_days = $3 WHERE report_id = $4 RETURNING *;

-- name: ListReportedContributors :many
SELECT 
    r.report_id,
    r.reported_by,
    r.target_user_id,
    r.reason,
    r.status,
    r.reviewed_at,
    r.reviewedby,
    r.created_at,

    -- Reporter Details
    ru.user_id AS reported_by_id,
    ru.name AS reported_by_name,
    ru.username AS reported_by_username,

    -- Target Contributor Details
    tu.user_id AS target_user_id,
    tu.name AS target_name,
    tu.username AS target_username,

    -- Post Details
    p.post_id AS post_id,
    p.slug AS target_post_slug,

    -- Comment Details
    c.comment_id AS comment_id,
    c.content AS target_comment,

    -- Moderator Details
    m.moderator_id AS reviewer_moderator_id,
    m.name AS reviewer_name

FROM reports r
LEFT JOIN users ru ON r.reported_by = ru.user_id  
LEFT JOIN users tu ON r.target_user_id = tu.user_id  
INNER JOIN contributors ct ON tu.user_id = ct.user_id
LEFT JOIN posts p ON r.target_post_id = p.post_id
LEFT JOIN comments c ON r.target_comment_id = c.comment_id  
LEFT JOIN moderators m ON r.reviewedby = m.moderator_id
WHERE r.target_user_id IS NOT NULL
ORDER BY r.created_at DESC;

-- name: ListReportedUsers :many
SELECT 
    r.report_id,
    r.reported_by,
    r.target_user_id,
    r.reason,
    r.status,
    r.reviewed_at,
    r.reviewedby,
    r.created_at,

    -- Reporter Details
    ru.user_id AS reported_by_id,
    ru.name AS reported_by_name,
    ru.username AS reported_by_username,

    -- Target User Details (Non-Contributors Only)
    tu.user_id AS target_user_id,
    tu.name AS target_name,
    tu.username AS target_username,

    -- Post Details
    p.post_id AS post_id,
    p.slug AS target_post_slug,

    -- Comment Details
    c.comment_id AS comment_id,
    c.content AS target_comment,

    -- Moderator Details
    m.moderator_id AS reviewer_moderator_id,
    m.name AS reviewer_name

FROM reports r
LEFT JOIN users ru ON r.reported_by = ru.user_id  
LEFT JOIN users tu ON r.target_user_id = tu.user_id  
LEFT JOIN contributors ct ON tu.user_id = ct.user_id  -- Check if the user is a contributor
LEFT JOIN posts p ON r.target_post_id = p.post_id
LEFT JOIN comments c ON r.target_comment_id = c.comment_id  
LEFT JOIN moderators m ON r.reviewedby = m.moderator_id
WHERE r.target_user_id IS NOT NULL
AND ct.user_id IS NULL  -- Ensures reported user is NOT a contributor
ORDER BY r.created_at DESC;

-- name: GetResolvedReportsWithSuspensionByUserId :many
SELECT DISTINCT ON (r.report_id)
    r.report_id,
    r.reported_by,
    r.target_user_id,
    r.reason,
    r.status AS report_status,
    r.reviewed_at,
    r.reviewedby,
    r.created_at,
    r.suspend_days,
    -- Target User Details
    tu.user_id AS target_user_id,
    tu.name AS target_name,
    tu.username AS target_username,
    -- Post Details
    p.post_id AS post_id,
    p.slug AS target_post_slug,
    -- Comment Details
    c.comment_id AS comment_id,
    c.content AS target_comment,
    -- Most recent appeal details
    a.appeal_id,
    a.status AS appeal_status,
    a.created_at AS appeal_created_at
FROM reports r
LEFT JOIN users tu ON r.target_user_id = tu.user_id  
LEFT JOIN posts p ON r.target_post_id = p.post_id
LEFT JOIN comments c ON r.target_comment_id = c.comment_id  
LEFT JOIN appeals a ON r.report_id = a.target_report_id AND a.status = 'dismissed'
WHERE r.status = 'resolved'
AND r.suspend_days IS NOT NULL
AND r.target_user_id = $1
AND NOT EXISTS (
    SELECT 1 FROM appeals 
    WHERE appeals.target_report_id = r.report_id 
    AND appeals.status = 'pending'
)
ORDER BY r.report_id, a.created_at DESC;

-- name: GetReportById :one
SELECT 
    r.report_id,
    r.reported_by,
    r.target_user_id,
    r.reason,
    r.status AS report_status,
    r.reviewed_at,
    r.reviewedby,
    r.created_at,
r.suspend_days
FROM reports r
WHERE r.report_id = $1;