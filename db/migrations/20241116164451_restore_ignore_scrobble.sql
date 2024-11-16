-- +goose Up
-- +goose StatementBegin
alter table media_file
	add column ignore_scrobble bool default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
