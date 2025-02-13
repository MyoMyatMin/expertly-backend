-- name: CreateAppeal :one
INSERT INTO appeals(appeal_id, appealed_by, target_report_id, reason) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateAppealStatus :one
UPDATE appeals SET status = $1, reviewedby = $2, reviewed_at = CURRENT_TIMESTAMP WHERE appeal_id = $3 RETURNING *;

-- name: ListAllAppealDetails :many
SELECT 
    a.appeal_id as appeal_id,
    a.appealed_by as appealed_by,
    a.target_report_id as target_report_id,
    a.reason as reason,
    a.status as status,
    a.reviewed_at as reviewed_at,
    a.reviewedby as reviewedby,
    a.created_at as created_at,

    u.user_id as appealed_by_id,
    u.name as appealed_by_name,

    r.report_id as target_report_id,
    r.reason as target_report_reason,
    r.target_post_id as target_post_id,
    r.target_user_id as target_user_id,
    r.target_comment_id as target_comment_id

FROM appeals a
LEFT JOIN users u ON a.appealed_by = u.user_id
LEFT JOIN reports r ON a.target_report_id = r.report_id
ORDER BY a.created_at DESC;
-- name: ListAppealsByContributors :many
SELECT 
    a.appeal_id as appeal_id,
    a.appealed_by as appealed_by,
    a.target_report_id as target_report_id,
    a.reason as reason,
    a.status as status,
    a.reviewed_at as reviewed_at,
    a.reviewedby as reviewedby,
    a.created_at as created_at,

    u.user_id as appealed_by_id,
    u.name as appealed_by_name,
    u.username as appealed_by_username,

    r.report_id as target_report_id,
    r.reason as target_report_reason,
    r.target_post_id as target_post_id,
    r.target_user_id as target_user_id,
    r.target_comment_id as target_comment_id,

    -- Moderator Details
    m.moderator_id AS reviewer_moderator_id,
    m.name AS reviewer_name

FROM appeals a
LEFT JOIN users u ON a.appealed_by = u.user_id
LEFT JOIN reports r ON a.target_report_id = r.report_id
INNER JOIN contributors c ON u.user_id = c.user_id  -- Check if user is a contributor
LEFT JOIN moderators m ON a.reviewedby = m.moderator_id
ORDER BY a.created_at DESC;

-- name: ListAppealsByUsers :many
SELECT 
    a.appeal_id as appeal_id,
    a.appealed_by as appealed_by,
    a.target_report_id as target_report_id,
    a.reason as reason,
    a.status as status,
    a.reviewed_at as reviewed_at,
    a.reviewedby as reviewedby,
    a.created_at as created_at,

    u.user_id as appealed_by_id,
    u.name as appealed_by_name,
    u.username as appealed_by_username,

    r.report_id as target_report_id,
    r.reason as target_report_reason,
    r.target_post_id as target_post_id,
    r.target_user_id as target_user_id,
    r.target_comment_id as target_comment_id,

    -- Moderator Details
    m.moderator_id AS reviewer_moderator_id,
    m.name AS reviewer_name

FROM appeals a
LEFT JOIN users u ON a.appealed_by = u.user_id
LEFT JOIN reports r ON a.target_report_id = r.report_id
LEFT JOIN contributors c ON u.user_id = c.user_id  -- Check if user is a contributor
LEFT JOIN moderators m ON a.reviewedby = m.moderator_id
WHERE c.user_id IS NULL  -- Ensures appealing user is NOT a contributor
ORDER BY a.created_at DESC;


-- name: GetAppealById :one
SELECT 
    -- Appeal fields
    a.appeal_id,
    a.appealed_by,
    a.target_report_id,
    a.reason as appeal_reason,
    a.status as appeal_status,
    a.reviewed_at,
    a.reviewedby,
    a.created_at,

    -- Appealing user fields
    u.user_id as appealed_by_id,
    u.name as appealed_by_name,

    -- Report fields
    r.report_id as target_report_id,
    r.reason as target_report_reason,
    r.target_post_id,
    r.target_user_id,
    r.target_comment_id,

    -- Added: Comment content
    c.content as comment_content,

    -- Added: Post slug
    p.slug as post_slug,

    -- Added: Target user details
    tu.suspended_until as target_user_suspended_until,
    tu.username as target_user_username,

    m.name as reviewer_name
FROM appeals a
LEFT JOIN users u ON a.appealed_by = u.user_id
LEFT JOIN reports r ON a.target_report_id = r.report_id
LEFT JOIN comments c ON r.target_comment_id = c.comment_id
LEFT JOIN posts p ON r.target_post_id = p.post_id
LEFT JOIN users tu ON r.target_user_id = tu.user_id
LEFT JOIN moderators m ON a.reviewedby = m.moderator_id
WHERE a.appeal_id = $1;

