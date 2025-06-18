-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- Comment out as these aren't used anymore
-- ALTER TABLE media_file
    -- ADD COLUMN explicit bool default false;

-- ALTER TABLE album
    -- ADD COLUMN explicit bool default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
