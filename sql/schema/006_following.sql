-- +goose Up
CREATE TABLE following (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    followee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    followed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    PRIMARY KEY (follower_id, followee_id) 
);

-- +goose Down
DROP TABLE following;