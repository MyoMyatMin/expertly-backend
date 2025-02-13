-- name: CreateReport :one
INSERT INTO reports (report_id, reported_by, target_post_id, target_user_id, target_comment_id, reason) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateReportStatus :one
UPDATE reports SET status = $1, reviewedby = $2, reviewed_at = CURRENT_TIMESTAMP WHERE report_id = $3 RETURNING *;

-- name: ListAllReportDetails :many
SELECT 
    r.report_id AS report_id,
    r.reported_by AS reported_by,
    r.target_post_id AS reported_target_post_id, -- Renamed to avoid duplication
    r.target_user_id AS target_user_id,
    r.target_comment_id AS reported_target_comment_id, -- Renamed to avoid duplication
    r.reason AS reason,
    r.status AS status,
    r.reviewed_at AS reviewed_at,
    r.reviewedby AS reviewedby,
    r.created_at AS created_at,

    -- Reporter Details
    ru.user_id AS reported_by_id,
    ru.name AS reported_by_name,
    ru.username AS reported_by_username,

    -- Target User Details
    tu.user_id AS target_user_id,
    tu.name AS target_name,
    tu.username AS target_username,

    -- Post Details
    p.post_id AS post_id,
    p.slug AS target_post_slug,

    -- Comment Details
    c.comment_id AS comment_id,
    c.content AS target_comment

FROM reports r
LEFT JOIN users ru ON r.reported_by = ru.user_id  -- Reporter (who reported)
LEFT JOIN users tu ON r.target_user_id = tu.user_id  -- Target User (who got reported)
LEFT JOIN posts p ON r.target_post_id = p.post_id
LEFT JOIN comments c ON r.target_comment_id = c.comment_id  
ORDER BY r.created_at DESC;

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