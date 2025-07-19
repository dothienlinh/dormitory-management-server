package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateContractStatusEnum, downUpdateContractStatusEnum)
}

func upUpdateContractStatusEnum(ctx context.Context, tx *sql.Tx) error {
	// Directly create new enum and update column
	query := `
	-- Create new enum with correct values
	CREATE TYPE contract_status_new AS ENUM ('active', 'inactive', 'cancelled');
	
	-- Update table to use new enum (convert through text to avoid issues)
	ALTER TABLE contracts ALTER COLUMN status TYPE contract_status_new USING (
		CASE 
			WHEN status::text = 'pending' THEN 'inactive'::contract_status_new
			WHEN status::text = 'expired' THEN 'cancelled'::contract_status_new
			WHEN status::text = 'active' THEN 'active'::contract_status_new
			WHEN status::text = 'cancelled' THEN 'cancelled'::contract_status_new
			ELSE 'inactive'::contract_status_new
		END
	);
	
	-- Drop old enum and rename new one
	DROP TYPE contract_status_enum;
	ALTER TYPE contract_status_new RENAME TO contract_status_enum;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downUpdateContractStatusEnum(ctx context.Context, tx *sql.Tx) error {
	// Revert to original enum values
	query := `
	-- Create original enum
	CREATE TYPE contract_status_old AS ENUM ('pending', 'active', 'expired', 'cancelled');
	
	-- Update table to use old enum with proper conversion
	ALTER TABLE contracts ALTER COLUMN status TYPE contract_status_old USING (
		CASE 
			WHEN status::text = 'inactive' THEN 'pending'::contract_status_old
			WHEN status::text = 'active' THEN 'active'::contract_status_old
			WHEN status::text = 'cancelled' THEN 'cancelled'::contract_status_old
			ELSE 'pending'::contract_status_old
		END
	);
	
	-- Drop new enum and rename old one
	DROP TYPE contract_status_enum;
	ALTER TYPE contract_status_old RENAME TO contract_status_enum;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
