package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddPlaylistExternalInfo, downAddPlaylistExternalInfo)
}

func upAddPlaylistExternalInfo(tx *sql.Tx) error {
	// Note: Ideally, we would also change the type of "comment" to be longer than 255
	// characters, but since this is Sqlite, the length doesn't matter
	_, err := tx.Exec(`
alter table playlist
	add external_agent varchar default '' not null;
alter table playlist 
	add external_id varchar default '' not null;
alter table playlist
	add external_url varchar default '' not null;
	`)
	return err
}

func downAddPlaylistExternalInfo(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
