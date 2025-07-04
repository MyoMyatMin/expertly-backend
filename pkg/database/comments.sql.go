// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: comments.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createComment = `-- name: CreateComment :one
INSERT INTO comments (
    comment_id,
    post_id,
    user_id,
    parent_comment_id,
    content
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING comment_id, post_id, user_id, parent_comment_id, content, created_at, updated_at
`

type CreateCommentParams struct {
	CommentID       uuid.UUID
	PostID          uuid.UUID
	UserID          uuid.UUID
	ParentCommentID uuid.NullUUID
	Content         string
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment,
		arg.CommentID,
		arg.PostID,
		arg.UserID,
		arg.ParentCommentID,
		arg.Content,
	)
	var i Comment
	err := row.Scan(
		&i.CommentID,
		&i.PostID,
		&i.UserID,
		&i.ParentCommentID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :exec
DELETE FROM comments
WHERE comment_id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteComment, commentID)
	return err
}

const getCommentByID = `-- name: GetCommentByID :one
SELECT 
    comment_id, 
    post_id, 
    user_id, 
    parent_comment_id, 
    content, 
    created_at, 
    updated_at
FROM comments
WHERE comment_id = $1
`

func (q *Queries) GetCommentByID(ctx context.Context, commentID uuid.UUID) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getCommentByID, commentID)
	var i Comment
	err := row.Scan(
		&i.CommentID,
		&i.PostID,
		&i.UserID,
		&i.ParentCommentID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCommentsByPost = `-- name: GetCommentsByPost :many
SELECT 
    c.comment_id, 
    c.post_id, 
    c.user_id, 
    c.parent_comment_id, 
    c.content, 
    c.created_at, 
    c.updated_at,
    u.username,
    u.name
FROM comments c
JOIN users u ON c.user_id = u.user_id
WHERE c.post_id = $1
ORDER BY c.created_at ASC
`

type GetCommentsByPostRow struct {
	CommentID       uuid.UUID
	PostID          uuid.UUID
	UserID          uuid.UUID
	ParentCommentID uuid.NullUUID
	Content         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Username        string
	Name            string
}

func (q *Queries) GetCommentsByPost(ctx context.Context, postID uuid.UUID) ([]GetCommentsByPostRow, error) {
	rows, err := q.db.QueryContext(ctx, getCommentsByPost, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsByPostRow
	for rows.Next() {
		var i GetCommentsByPostRow
		if err := rows.Scan(
			&i.CommentID,
			&i.PostID,
			&i.UserID,
			&i.ParentCommentID,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.Name,
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

const updateComment = `-- name: UpdateComment :one
UPDATE comments 
SET 
    content = $2, 
    updated_at = CURRENT_TIMESTAMP
WHERE comment_id = $1
RETURNING comment_id, post_id, user_id, parent_comment_id, content, created_at, updated_at
`

type UpdateCommentParams struct {
	CommentID uuid.UUID
	Content   string
}

func (q *Queries) UpdateComment(ctx context.Context, arg UpdateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, updateComment, arg.CommentID, arg.Content)
	var i Comment
	err := row.Scan(
		&i.CommentID,
		&i.PostID,
		&i.UserID,
		&i.ParentCommentID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
