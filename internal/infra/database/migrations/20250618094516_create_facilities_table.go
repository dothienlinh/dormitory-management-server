package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateFacilitiesTable, downCreateFacilitiesTable)
}

func upCreateFacilitiesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE facilities (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downCreateFacilitiesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE facilities;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
