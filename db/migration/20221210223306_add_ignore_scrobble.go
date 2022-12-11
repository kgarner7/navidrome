package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddIgnoreScrobble, downAddIgnoreScrobble)
}

func upAddIgnoreScrobble(tx *sql.Tx) error {
	_, err := tx.Exec(`
alter table media_file
	add column ignore_scrobble bool default false;
`)
	if err != nil {
		return err
	}
	notice(tx, "A full rescan needs to be performed to ignore scrobble tabs")
	return forceFullRescan(tx)
}

func downAddIgnoreScrobble(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
