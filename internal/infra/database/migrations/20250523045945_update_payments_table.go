package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdatePaymentsTable, downUpdatePaymentsTable)
}

func upUpdatePaymentsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TYPE payment_method RENAME TO payment_channel_enum;

	ALTER TABLE payments RENAME COLUMN payment_method TO payment_channel;

	CREATE TYPE payment_method_enum AS ENUM ('VIETQR', 'MOMO', 'VNPAY');
	
	ALTER TABLE payments
	ADD COLUMN payment_method payment_method_enum NOT NULL;

	ALTER TYPE user_role RENAME TO user_role_enum;

	ALTER TYPE user_gender RENAME TO user_gender_enum;

	ALTER TYPE user_status RENAME TO user_status_enum;

	ALTER TYPE room_status RENAME TO room_status_enum;

	ALTER TYPE contract_status RENAME TO contract_status_enum;

	ALTER TYPE payment_status RENAME TO payment_status_enum;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downUpdatePaymentsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	ALTER TYPE payment_status_enum RENAME TO payment_status;

	ALTER TYPE contract_status_enum RENAME TO contract_status;

	ALTER TYPE room_status_enum RENAME TO room_status;

	ALTER TYPE user_status_enum RENAME TO user_status;

	ALTER TYPE user_gender_enum RENAME TO user_gender;

	ALTER TABLE payments DROP COLUMN payment_method;

	ALTER TYPE user_role_enum RENAME TO user_role;

	DROP TYPE payment_method_enum;

	ALTER TABLE payments RENAME COLUMN payment_channel TO payment_method;

	ALTER TYPE payment_channel_enum RENAME TO payment_method;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
