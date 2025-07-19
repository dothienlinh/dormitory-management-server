package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateContractTermsTable, downCreateContractTermsTable)
}

func upCreateContractTermsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE contract_terms (
		id SERIAL PRIMARY KEY,
		contract_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		"order" INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		CONSTRAINT fk_contract_terms_contract_id FOREIGN KEY (contract_id) REFERENCES contracts(id) ON DELETE CASCADE
	);

	CREATE INDEX idx_contract_terms_contract_id ON contract_terms (contract_id);
	CREATE INDEX idx_contract_terms_order ON contract_terms ("order");
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downCreateContractTermsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE IF EXISTS contract_terms;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
