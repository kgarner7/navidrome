-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE media_file
    ADD COLUMN explicit bool default false;

ALTER TABLE album
    ADD COLUMN explicit bool default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
