-- +goose Up
CREATE TABLE reports (
    report_id UUID PRIMARY KEY,
    reported_by UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    target_post_id UUID REFERENCES posts(post_id) ON DELETE CASCADE,
    target_user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    target_comment_id UUID REFERENCES comments(comment_id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'resolved', 'dismissed')) DEFAULT 'pending',
    reviewed_at TIMESTAMP,
    reviewedby UUID REFERENCES moderators(moderator_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE reports;