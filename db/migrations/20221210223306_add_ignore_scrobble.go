package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddIgnoreScrobble, downAddIgnoreScrobble)
}

func upAddIgnoreScrobble(_ context.Context, tx *sql.Tx) error {
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

func downAddIgnoreScrobble(_ context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
