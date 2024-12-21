-- +goose Up
CREATE TABLE posts(
    id UUID primary key,
    user_id UUID references users(id) ON DELETE CASCADE ,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'published'))
);

-- +goose Down
DROP TABLE posts;
