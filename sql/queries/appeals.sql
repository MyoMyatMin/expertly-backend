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


-- name: GetAppealById :one
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
WHERE a.appeal_id = $1;

