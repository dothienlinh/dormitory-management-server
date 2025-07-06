package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddStatusAccountFieldForUserTable, downAddStatusAccountFieldForUserTable)
}

func upAddStatusAccountFieldForUserTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TYPE status_account_enum AS ENUM ('pending', 'approved', 'rejected', 'banned');

	ALTER TABLE users
	ADD COLUMN status_account status_account_enum NOT NULL DEFAULT 'pending';
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downAddStatusAccountFieldForUserTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE users
	DROP COLUMN status_account;

	DROP TYPE IF EXISTS status_account_enum;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
