package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddAddressField, downAddAddressField)
}

func upAddAddressField(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE users
	ADD COLUMN address VARCHAR(255) NULL DEFAULT NULL;
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downAddAddressField(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE users
	DROP COLUMN address;
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}
