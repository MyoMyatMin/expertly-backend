-- +goose Up
CREATE TABLE saved_posts(
    user_id UUID NOT NULL REFERENCES users(user_id),
    post_id UUID NOT NULL REFERENCES posts(post_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id)
);

-- +goose Down
DROP TABLE saved_posts;