package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddFilePriceForContractsTable, downAddFilePriceForContractsTable)
}

func upAddFilePriceForContractsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE contracts
	ADD COLUMN price FLOAT NOT NULL DEFAULT 0,
	ADD COLUMN description TEXT NOT NULL DEFAULT '';
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downAddFilePriceForContractsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE contracts
	DROP COLUMN price,
	DROP COLUMN description
	;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
