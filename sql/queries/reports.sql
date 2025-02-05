-- name: CreateReport :one
INSERT INTO reports (report_id, reported_by, target_post_id, target_user_id, target_comment_id, reason) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateReportStatus :one
UPDATE reports SET status = $1, reviewed_at = $2, reviewedby = $3, updated_at = CURRENT_TIMESTAMP WHERE report_id = $4 RETURNING *;

-- name: ListAllReportDetails :many
SELECT 
    r.report_id as report_id,
    r.reported_by as reported_by,
    r.target_post_id as target_post_id,
    r.target_user_id as target_user_id,
    r.target_comment_id as target_comment_id,
    r.reason as reason,
    r.status as status,
    r.reviewed_at as reviewed_at,
    r.reviewedby as reviewedby,
    r.created_at as created_at,

    u.user_id as reported_by_id,
    u.name as reported_by_name,
    u.username as reported_by_username,

    p.post_id as target_post_id,
    p.slug as target_post_slug,

    c.comment_id as target_comment_id,
    c.content as target_comment

FROM reports r
LEFT JOIN users u ON r.reported_by = u.id
LEFT JOIN posts p ON r.target_post_id = p.post_id
LEFT JOIN comments c ON r.target_comment_id = c.id
ORDER BY r.created_at DESC;

