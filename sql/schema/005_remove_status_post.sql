-- +goose Up
ALTER TABLE posts
DROP COLUMN status;

-- +goose Down
ALTER TABLE posts
ADD COLUMN status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'published'));
