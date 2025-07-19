package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreatePaymentHistoriesTable, downCreatePaymentHistoriesTable)
}

func upCreatePaymentHistoriesTable(ctx context.Context, tx *sql.Tx) error {
	// Create payment_status enum if it doesn't exist
	createStatusEnum := `
	DO $$ BEGIN
		CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'overdue', 'cancelled');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;
	`
	if _, err := tx.ExecContext(ctx, createStatusEnum); err != nil {
		return err
	}

	// Create payment_method enum if it doesn't exist
	createMethodEnum := `
	DO $$ BEGIN
		CREATE TYPE payment_method AS ENUM ('cash', 'bank_transfer', 'card', 'mobile_payment');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;
	`
	if _, err := tx.ExecContext(ctx, createMethodEnum); err != nil {
		return err
	}

	// Create payment_histories table
	query := `
	CREATE TABLE payment_histories (
		id SERIAL PRIMARY KEY,
		contract_id INTEGER NOT NULL,
		period VARCHAR(255) NOT NULL,
		amount FLOAT NOT NULL,
		status payment_status NOT NULL DEFAULT 'pending',
		payment_date TIMESTAMPTZ,
		due_date TIMESTAMPTZ NOT NULL,
		method payment_method,
		description TEXT,
		receipt_url TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		CONSTRAINT fk_payment_histories_contract_id FOREIGN KEY (contract_id) REFERENCES contracts(id) ON DELETE CASCADE
	);

	CREATE INDEX idx_payment_histories_contract_id ON payment_histories (contract_id);
	CREATE INDEX idx_payment_histories_status ON payment_histories (status);
	CREATE INDEX idx_payment_histories_due_date ON payment_histories (due_date);
	CREATE INDEX idx_payment_histories_payment_date ON payment_histories (payment_date);
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downCreatePaymentHistoriesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE IF EXISTS payment_histories;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	// Note: We don't drop the enums as they might be used by other tables
	// If you need to drop them, uncomment the following:
	/*
		dropEnums := `
		DROP TYPE IF EXISTS payment_status;
		DROP TYPE IF EXISTS payment_method;
		`
		if _, err := tx.ExecContext(ctx, dropEnums); err != nil {
			return err
		}
	*/

	return nil
}
