-- +goose Up
-- +goose StatementBegin
ALTER TABLE comments
ALTER COLUMN parent_id DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE comments 
ALTER COLUMN parent_id SET NOT NULL;
-- +goose StatementEnd
