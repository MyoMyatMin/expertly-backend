-- +goose Up
CREATE TABLE Contributors (
    user_id UUID PRIMARY KEY REFERENCES Users(user_id) ON DELETE CASCADE,
    expertise_fields TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE contributors;