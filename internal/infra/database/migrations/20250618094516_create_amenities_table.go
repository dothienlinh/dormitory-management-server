package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateAmenitiesTable, downCreateAmenitiesTable)
}

func upCreateAmenitiesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE amenities (
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

func downCreateAmenitiesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE amenities;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
