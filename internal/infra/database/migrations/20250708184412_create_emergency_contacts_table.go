package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateEmergencyContactsTable, downCreateEmergencyContactsTable)
}

func upCreateEmergencyContactsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE emergency_contacts (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		name VARCHAR(255) NOT NULL,
		phone VARCHAR(15) NOT NULL,
		relationship VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

	)
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateEmergencyContactsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE emergency_contacts
	`

	_, err := tx.ExecContext(ctx, query)
	return err
}
