package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upBillsTable, downBillsTable)
}

func upBillsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TYPE bill_status_enum AS ENUM ('PENDING', 'PAID', 'OVERDUE');

	CREATE TABLE IF NOT EXISTS bills (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		payment_id INT,
		amount DECIMAL(10, 2) NOT NULL,
		status bill_status_enum NOT NULL DEFAULT 'PENDING',
		description TEXT,
		due_date TIMESTAMP NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE SET NULL
	);
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func downBillsTable(ctx context.Context, tx *sql.Tx) error {
	dropTableQuery := `
	DROP TABLE IF EXISTS bills;

	DROP TYPE IF EXISTS bill_status_enum;
	`
	if _, err := tx.ExecContext(ctx, dropTableQuery); err != nil {
		return err
	}

	return nil
}
