package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateMaintenanceHistoriesTable, downCreateMaintenanceHistoriesTable)
}

func upCreateMaintenanceHistoriesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE maintenance_histories (
		id SERIAL PRIMARY KEY,
		description TEXT NOT NULL,
		room_id BIGINT NOT NULL,
		maintenance_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		cost FLOAT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
	);
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downCreateMaintenanceHistoriesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE maintenance_histories;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
