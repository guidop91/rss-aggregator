-- +goose Up
CREATE TABLE users (
  id UUID NOT NULL PRIMARY KEY,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  name text NOT NULL
);

-- +goose Down
DROP TABLE users;
