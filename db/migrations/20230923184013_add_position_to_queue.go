package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddPositionToQueue, downAddPositionToQueue)
}

func upAddPositionToQueue(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`alter table playqueue add queue_index integer default 0;`)
	return err
}

func downAddPositionToQueue(_ context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
