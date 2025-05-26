package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upPaymentTable, downPaymentTable)
}

func upPaymentTable(ctx context.Context, tx *sql.Tx) error {
	createEnum := `
	CREATE TYPE payment_status AS ENUM ('pending', 'success', 'failed');
	CREATE TYPE payment_method AS ENUM ('credit_card', 'debit_card', 'bank_transfer');
	`
	if _, err := tx.ExecContext(ctx, createEnum); err != nil {
		return err
	}

	createTable := `
	CREATE TABLE payments (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT NOT NULL,
		currency VARCHAR(3) NOT NULL,
		payment_status payment_status NOT NULL,
		payment_method payment_method NOT NULL,
		transaction_id VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		CONSTRAINT fk_payments_user_id FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE INDEX idx_payments_transaction_id ON payments (transaction_id);
	`
	if _, err := tx.ExecContext(ctx, createTable); err != nil {
		return err
	}

	removeContractIDFromUsers := `
	ALTER TABLE users
	DROP CONSTRAINT fk_contracts_user;

	ALTER TABLE users
	DROP COLUMN contract_id;
	`
	if _, err := tx.ExecContext(ctx, removeContractIDFromUsers); err != nil {
		return err
	}

	return nil
}

func downPaymentTable(ctx context.Context, tx *sql.Tx) error {
	dropTable := `
	DROP TABLE payments;
	`

	_, err := tx.ExecContext(ctx, dropTable)
	if err != nil {
		return err
	}

	dropEnum := `
	DROP TYPE payment_status;
	DROP TYPE payment_method;
	`
	if _, err := tx.ExecContext(ctx, dropEnum); err != nil {
		return err
	}

	return nil
}
