package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateContractsTable, downUpdateContractsTable)
}

func upUpdateContractsTable(ctx context.Context, tx *sql.Tx) error {
	// Create payment_cycle enum
	createEnum := `
	CREATE TYPE payment_cycle AS ENUM ('monthly', 'quarterly', 'yearly');
	`
	if _, err := tx.ExecContext(ctx, createEnum); err != nil {
		return err
	}

	// Add new fields to contracts table
	query := `
	ALTER TABLE contracts
	ADD COLUMN signed_date TIMESTAMPTZ,
	ADD COLUMN deposit_amount FLOAT,
	ADD COLUMN payment_cycle payment_cycle DEFAULT 'monthly';
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	// Update existing contracts to have default status if needed
	updateQuery := `
	UPDATE contracts 
	SET status = 'active' 
	WHERE status IS NULL;
	`
	if _, err := tx.ExecContext(ctx, updateQuery); err != nil {
		return err
	}

	return nil
}

func downUpdateContractsTable(ctx context.Context, tx *sql.Tx) error {
	// Remove added fields
	query := `
	ALTER TABLE contracts
	DROP COLUMN signed_date,
	DROP COLUMN deposit_amount,
	DROP COLUMN payment_cycle;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	// Drop payment_cycle enum
	dropEnum := `
	DROP TYPE payment_cycle;
	`
	if _, err := tx.ExecContext(ctx, dropEnum); err != nil {
		return err
	}

	return nil
}
