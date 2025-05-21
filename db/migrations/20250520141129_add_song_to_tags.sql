-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE tag
    ADD COLUMN song bool default false;
CREATE INDEX IF NOT EXISTS tag_song
    ON tag (song);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
