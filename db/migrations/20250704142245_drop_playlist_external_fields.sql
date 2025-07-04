-- +goose Up
-- +goose StatementBegin
ALTER TABLE playlist DROP COLUMN external_agent;
ALTER TABLE playlist DROP COLUMN external_id;
ALTER TABLE playlist DROP COLUMN external_url;
ALTER TABLE playlist DROP COLUMN external_sync;
ALTER TABLE playlist DROP COLUMN external_syncable;
ALTER TABLE playlist DROP COLUMN external_recommended;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
