-- +goose Up
CREATE TABLE upvotes (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id) -- Prevent duplicate upvotes
);


-- +goose Down
DROP TABLE upvotes;