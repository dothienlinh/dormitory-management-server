package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdatePaymentTable, downUpdatePaymentTable)
}

func upUpdatePaymentTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TABLE payments
	DROP COLUMN payment_status;

	DROP TYPE payment_status_enum;

	ALTER TABLE payments
	ALTER COLUMN transaction_id DROP NOT NULL;

	ALTER TABLE payments
	ADD COLUMN bin text NOT NULL DEFAULT '',
	ADD COLUMN account_number text NOT NULL DEFAULT '',
	ADD COLUMN account_name text NOT NULL DEFAULT '',
	ADD COLUMN description text NOT NULL DEFAULT '',
	ADD COLUMN order_code bigint NOT NULL DEFAULT 0,
	ADD COLUMN payment_link_id text NOT NULL DEFAULT '',
	ADD COLUMN status text NOT NULL DEFAULT 'PENDING',
	ADD COLUMN expired_at int DEFAULT NULL;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downUpdatePaymentTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TYPE payment_status_enum AS ENUM ('pending', 'success', 'failed');

	ALTER TABLE payments
	ADD COLUMN payment_status payment_status_enum NOT NULL DEFAULT 'pending';

	ALTER TABLE payments
	ALTER COLUMN transaction_id SET NOT NULL;
	
	ALTER TABLE payments
	DROP COLUMN bin,
	DROP COLUMN account_number,
	DROP COLUMN account_name,
	DROP COLUMN description,
	DROP COLUMN order_code,
	DROP COLUMN payment_link_id,
	DROP COLUMN status,
	DROP COLUMN expired_at;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
