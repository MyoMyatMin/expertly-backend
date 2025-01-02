-- +goose Up

-- Drop the `status` column in `posts`
ALTER TABLE posts
DROP COLUMN status;

-- Add `username` to `users` table
ALTER TABLE users
ADD COLUMN username TEXT NOT NULL UNIQUE;

-- Add `slug` to `posts` table
ALTER TABLE posts
ADD COLUMN slug TEXT NOT NULL UNIQUE;

-- +goose Down

-- Add the `status` column back to `posts`
ALTER TABLE posts
ADD COLUMN status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'published'));

-- Remove `username` from `users` table
ALTER TABLE users
DROP COLUMN username;

-- Remove `slug` from `posts` table
ALTER TABLE posts
DROP COLUMN slug;
