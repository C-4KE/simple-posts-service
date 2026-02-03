-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXIST ltree;

CREATE TABLE IF NOT EXISTS Comments (
    comment_id BIGSERIAL PRIMARY KEY,
    author_id UUID NOT NULL,
    post_id BIGINT NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    parent_id BIGINT NOT NULL REFERENCES comments(comment_id) ON DELETE CASCADE,
    text VARCHAR(2000) NOT NULL,
    create_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    path LTREE UNIQUE NOT NULL,
    replies_level INTEGER NOT NULL
);

CREATE INDEX path_gist_idx ON comments USING GIST (path);

CREATE INDEX create_date_idx ON comments(create_date);

CREATE INDEX post_id_idx ON comments(post_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Comments;
-- +goose StatementEnd
