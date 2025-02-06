-- +goose Up

CREATE TABLE contributor_applications(
    contri_app_id UUID primary key,
    user_id UUID NOT NULL REFERENCES users(user_id),
    expertise_proofs TEXT[] NOT NULL,
    identity_proof TEXT NOT NULL,
    initial_submission TEXT NOT NULL,
    reviewed_by UUID REFERENCES moderators(moderator_id),
    status TEXT CHECK (status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP
);

-- +goose Down
DRoP TABLE contributor_applications;