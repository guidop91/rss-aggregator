-- +goose Up
CREATE TABLE posts (
  id UUID PRIMARY KEY,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  title text NOT NULL,
  description text,
  pubdate timestamp NOT NULL,
  url text NOT NULL UNIQUE,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
