package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddFieldIsVerifyInUsersTable, downAddFieldIsVerifyInUsersTable)
}

func upAddFieldIsVerifyInUsersTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE users
	ADD COLUMN is_verify BOOLEAN NOT NULL DEFAULT FALSE;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downAddFieldIsVerifyInUsersTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE users DROP COLUMN is_verify;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
