package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddStatusAccountFieldForUserTable, downAddStatusAccountFieldForUserTable)
}

func upAddStatusAccountFieldForUserTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TYPE status_account_enum AS ENUM ('pending', 'approved', 'rejected', 'banned');

	ALTER TABLE users
	ADD COLUMN status_account status_account_enum NOT NULL DEFAULT 'pending';

	ALTER TABLE users
	ALTER COLUMN student_code DROP NOT NULL;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downAddStatusAccountFieldForUserTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	UPDATE users 
	SET student_code = 'TEMP_' || EXTRACT(EPOCH FROM NOW())::bigint || '_' || id 
	WHERE student_code IS NULL;

	ALTER TABLE users
	ALTER COLUMN student_code SET NOT NULL;

	ALTER TABLE users
	DROP COLUMN status_account;

	DROP TYPE IF EXISTS status_account_enum;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
