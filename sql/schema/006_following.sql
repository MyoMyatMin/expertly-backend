-- +goose Up
CREATE TABLE following (
    follower_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, 
    following_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, 
    PRIMARY KEY (follower_id, following_id) 
);

-- +goose Down
DROP TABLE following;