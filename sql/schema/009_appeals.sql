-- +goose Up
CREATE TABLE appeals(
    appeal_id UUID PRIMARY KEY,
    appealed_by UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    target_report_id UUID NOT NULL REFERENCES reports(report_id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'resolved', 'dismissed')) DEFAULT 'pending',    
    reviewed_at TIMESTAMP,
    reviewedby UUID REFERENCES moderators(moderator_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE appeals;
