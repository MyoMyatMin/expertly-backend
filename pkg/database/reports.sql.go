// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: reports.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createReport = `-- name: CreateReport :one
INSERT INTO reports (report_id, reported_by, target_post_id, target_user_id, target_comment_id, reason) VALUES ($1, $2, $3, $4, $5, $6) RETURNING report_id, reported_by, target_post_id, target_user_id, target_comment_id, reason, status, reviewed_at, reviewedby, created_at, suspend_days
`

type CreateReportParams struct {
	ReportID        uuid.UUID
	ReportedBy      uuid.UUID
	TargetPostID    uuid.NullUUID
	TargetUserID    uuid.UUID
	TargetCommentID uuid.NullUUID
	Reason          string
}

func (q *Queries) CreateReport(ctx context.Context, arg CreateReportParams) (Report, error) {
	row := q.db.QueryRowContext(ctx, createReport,
		arg.ReportID,
		arg.ReportedBy,
		arg.TargetPostID,
		arg.TargetUserID,
		arg.TargetCommentID,
		arg.Reason,
	)
	var i Report
	err := row.Scan(
		&i.ReportID,
		&i.ReportedBy,
		&i.TargetPostID,
		&i.TargetUserID,
		&i.TargetCommentID,
		&i.Reason,
		&i.Status,
		&i.ReviewedAt,
		&i.Reviewedby,
		&i.CreatedAt,
		&i.SuspendDays,
	)
	return i, err
}

const getReportById = `-- name: GetReportById :one
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
WHERE r.report_id = $1
`

type GetReportByIdRow struct {
	ReportID     uuid.UUID
	ReportedBy   uuid.UUID
	TargetUserID uuid.UUID
	Reason       string
	ReportStatus sql.NullString
	ReviewedAt   sql.NullTime
	Reviewedby   uuid.NullUUID
	CreatedAt    sql.NullTime
	SuspendDays  sql.NullInt32
}

func (q *Queries) GetReportById(ctx context.Context, reportID uuid.UUID) (GetReportByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getReportById, reportID)
	var i GetReportByIdRow
	err := row.Scan(
		&i.ReportID,
		&i.ReportedBy,
		&i.TargetUserID,
		&i.Reason,
		&i.ReportStatus,
		&i.ReviewedAt,
		&i.Reviewedby,
		&i.CreatedAt,
		&i.SuspendDays,
	)
	return i, err
}

const getResolvedReportsWithSuspensionByUserId = `-- name: GetResolvedReportsWithSuspensionByUserId :many
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
ORDER BY r.report_id, a.created_at DESC
`

type GetResolvedReportsWithSuspensionByUserIdRow struct {
	ReportID        uuid.UUID
	ReportedBy      uuid.UUID
	TargetUserID    uuid.UUID
	Reason          string
	ReportStatus    sql.NullString
	ReviewedAt      sql.NullTime
	Reviewedby      uuid.NullUUID
	CreatedAt       sql.NullTime
	SuspendDays     sql.NullInt32
	TargetUserID_2  uuid.NullUUID
	TargetName      sql.NullString
	TargetUsername  sql.NullString
	PostID          uuid.NullUUID
	TargetPostSlug  sql.NullString
	CommentID       uuid.NullUUID
	TargetComment   sql.NullString
	AppealID        uuid.NullUUID
	AppealStatus    sql.NullString
	AppealCreatedAt sql.NullTime
}

func (q *Queries) GetResolvedReportsWithSuspensionByUserId(ctx context.Context, targetUserID uuid.UUID) ([]GetResolvedReportsWithSuspensionByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getResolvedReportsWithSuspensionByUserId, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetResolvedReportsWithSuspensionByUserIdRow
	for rows.Next() {
		var i GetResolvedReportsWithSuspensionByUserIdRow
		if err := rows.Scan(
			&i.ReportID,
			&i.ReportedBy,
			&i.TargetUserID,
			&i.Reason,
			&i.ReportStatus,
			&i.ReviewedAt,
			&i.Reviewedby,
			&i.CreatedAt,
			&i.SuspendDays,
			&i.TargetUserID_2,
			&i.TargetName,
			&i.TargetUsername,
			&i.PostID,
			&i.TargetPostSlug,
			&i.CommentID,
			&i.TargetComment,
			&i.AppealID,
			&i.AppealStatus,
			&i.AppealCreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listReportedContributors = `-- name: ListReportedContributors :many
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
ORDER BY r.created_at DESC
`

type ListReportedContributorsRow struct {
	ReportID            uuid.UUID
	ReportedBy          uuid.UUID
	TargetUserID        uuid.UUID
	Reason              string
	Status              sql.NullString
	ReviewedAt          sql.NullTime
	Reviewedby          uuid.NullUUID
	CreatedAt           sql.NullTime
	ReportedByID        uuid.NullUUID
	ReportedByName      sql.NullString
	ReportedByUsername  sql.NullString
	TargetUserID_2      uuid.NullUUID
	TargetName          sql.NullString
	TargetUsername      sql.NullString
	PostID              uuid.NullUUID
	TargetPostSlug      sql.NullString
	CommentID           uuid.NullUUID
	TargetComment       sql.NullString
	ReviewerModeratorID uuid.NullUUID
	ReviewerName        sql.NullString
}

func (q *Queries) ListReportedContributors(ctx context.Context) ([]ListReportedContributorsRow, error) {
	rows, err := q.db.QueryContext(ctx, listReportedContributors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListReportedContributorsRow
	for rows.Next() {
		var i ListReportedContributorsRow
		if err := rows.Scan(
			&i.ReportID,
			&i.ReportedBy,
			&i.TargetUserID,
			&i.Reason,
			&i.Status,
			&i.ReviewedAt,
			&i.Reviewedby,
			&i.CreatedAt,
			&i.ReportedByID,
			&i.ReportedByName,
			&i.ReportedByUsername,
			&i.TargetUserID_2,
			&i.TargetName,
			&i.TargetUsername,
			&i.PostID,
			&i.TargetPostSlug,
			&i.CommentID,
			&i.TargetComment,
			&i.ReviewerModeratorID,
			&i.ReviewerName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listReportedUsers = `-- name: ListReportedUsers :many
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
ORDER BY r.created_at DESC
`

type ListReportedUsersRow struct {
	ReportID            uuid.UUID
	ReportedBy          uuid.UUID
	TargetUserID        uuid.UUID
	Reason              string
	Status              sql.NullString
	ReviewedAt          sql.NullTime
	Reviewedby          uuid.NullUUID
	CreatedAt           sql.NullTime
	ReportedByID        uuid.NullUUID
	ReportedByName      sql.NullString
	ReportedByUsername  sql.NullString
	TargetUserID_2      uuid.NullUUID
	TargetName          sql.NullString
	TargetUsername      sql.NullString
	PostID              uuid.NullUUID
	TargetPostSlug      sql.NullString
	CommentID           uuid.NullUUID
	TargetComment       sql.NullString
	ReviewerModeratorID uuid.NullUUID
	ReviewerName        sql.NullString
}

func (q *Queries) ListReportedUsers(ctx context.Context) ([]ListReportedUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listReportedUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListReportedUsersRow
	for rows.Next() {
		var i ListReportedUsersRow
		if err := rows.Scan(
			&i.ReportID,
			&i.ReportedBy,
			&i.TargetUserID,
			&i.Reason,
			&i.Status,
			&i.ReviewedAt,
			&i.Reviewedby,
			&i.CreatedAt,
			&i.ReportedByID,
			&i.ReportedByName,
			&i.ReportedByUsername,
			&i.TargetUserID_2,
			&i.TargetName,
			&i.TargetUsername,
			&i.PostID,
			&i.TargetPostSlug,
			&i.CommentID,
			&i.TargetComment,
			&i.ReviewerModeratorID,
			&i.ReviewerName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateReportStatus = `-- name: UpdateReportStatus :one
UPDATE reports SET status = $1, reviewedby = $2, reviewed_at = CURRENT_TIMESTAMP, suspend_days = $3 WHERE report_id = $4 RETURNING report_id, reported_by, target_post_id, target_user_id, target_comment_id, reason, status, reviewed_at, reviewedby, created_at, suspend_days
`

type UpdateReportStatusParams struct {
	Status      sql.NullString
	Reviewedby  uuid.NullUUID
	SuspendDays sql.NullInt32
	ReportID    uuid.UUID
}

func (q *Queries) UpdateReportStatus(ctx context.Context, arg UpdateReportStatusParams) (Report, error) {
	row := q.db.QueryRowContext(ctx, updateReportStatus,
		arg.Status,
		arg.Reviewedby,
		arg.SuspendDays,
		arg.ReportID,
	)
	var i Report
	err := row.Scan(
		&i.ReportID,
		&i.ReportedBy,
		&i.TargetPostID,
		&i.TargetUserID,
		&i.TargetCommentID,
		&i.Reason,
		&i.Status,
		&i.ReviewedAt,
		&i.Reviewedby,
		&i.CreatedAt,
		&i.SuspendDays,
	)
	return i, err
}
