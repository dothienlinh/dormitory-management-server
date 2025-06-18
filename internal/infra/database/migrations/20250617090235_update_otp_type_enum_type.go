package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateOtpTypeEnumType, downUpdateOtpTypeEnumType)
}

func upUpdateOtpTypeEnumType(ctx context.Context, tx *sql.Tx) error {
	query := `
		ALTER TYPE otp_type_enum RENAME TO _otp_type_enum_old;
		
		CREATE TYPE otp_type_enum AS ENUM ('login', 'register', 'reset_password', 'verify_email', 'verify_phone', 'transaction', 'verify_account');
		
		ALTER TABLE otp_codes RENAME COLUMN otp_type TO _otp_type_old;
		
		ALTER TABLE otp_codes ADD COLUMN otp_type otp_type_enum;
		
		UPDATE otp_codes
		SET otp_type = CASE 
			WHEN _otp_type_old::text = 'login' THEN 'login'::otp_type_enum
			WHEN _otp_type_old::text = 'register' THEN 'register'::otp_type_enum
			WHEN _otp_type_old::text = 'reset_password' THEN 'reset_password'::otp_type_enum
			WHEN _otp_type_old::text = 'verify_email' THEN 'verify_email'::otp_type_enum
			WHEN _otp_type_old::text = 'verify_phone' THEN 'verify_phone'::otp_type_enum
			WHEN _otp_type_old::text = 'transaction' THEN 'transaction'::otp_type_enum
			WHEN _otp_type_old::text = 'verify_account' THEN 'verify_account'::otp_type_enum
			ELSE 'login'::otp_type_enum
		END;
		
		ALTER TABLE otp_codes ALTER COLUMN otp_type SET NOT NULL;
		
		ALTER TABLE otp_codes DROP COLUMN _otp_type_old;
		
		DROP TYPE _otp_type_enum_old;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downUpdateOtpTypeEnumType(ctx context.Context, tx *sql.Tx) error {
	query := `
		ALTER TYPE otp_type_enum RENAME TO _otp_type_enum_new;
		
		CREATE TYPE otp_type_enum AS ENUM ('login', 'register', 'reset_password', 'verify_email', 'verify_phone');
		
		ALTER TABLE otp_codes RENAME COLUMN otp_type TO _otp_type_new;
		
		ALTER TABLE otp_codes ADD COLUMN otp_type otp_type_enum;
		
		UPDATE otp_codes
		SET otp_type = CASE 
			WHEN _otp_type_new::text = 'login' THEN 'login'::otp_type_enum
			WHEN _otp_type_new::text = 'register' THEN 'register'::otp_type_enum
			WHEN _otp_type_new::text = 'reset_password' THEN 'reset_password'::otp_type_enum
			WHEN _otp_type_new::text = 'verify_email' THEN 'verify_email'::otp_type_enum
			WHEN _otp_type_new::text = 'verify_phone' THEN 'verify_phone'::otp_type_enum
			WHEN _otp_type_new::text = 'transaction' THEN 'login'::otp_type_enum
			WHEN _otp_type_new::text = 'verify_account' THEN 'verify_email'::otp_type_enum
			ELSE 'login'::otp_type_enum
		END;
		
		ALTER TABLE otp_codes ALTER COLUMN otp_type SET NOT NULL;
		
		ALTER TABLE otp_codes DROP COLUMN _otp_type_new;
		
		DROP TYPE _otp_type_enum_new;
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
