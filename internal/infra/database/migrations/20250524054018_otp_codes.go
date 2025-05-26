package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upOtpCodes, downOtpCodes)
}

func upOtpCodes(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TYPE otp_type_enum AS ENUM ('login', 'register', 'reset_password', 'verify_email', 'verify_phone', 'transaction');

	CREATE TYPE identifier_type_enum AS ENUM ('email', 'phone');

	CREATE TABLE otp_codes (
		id SERIAL PRIMARY KEY,
		user_id INT,
		identifier VARCHAR(255) NOT NULL,
		identifier_type identifier_type_enum NOT NULL,
		otp_code VARCHAR(10) NOT NULL,
		otp_type otp_type_enum NOT NULL,
		is_used BOOLEAN DEFAULT FALSE,
		expires_at TIMESTAMPTZ NOT NULL,
		verified_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		CONSTRAINT fk_opt_code_user_id FOREIGN KEY (user_id) REFERENCES users(id)
		);

		CREATE INDEX idx_otp_code ON otp_codes (otp_code);
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downOtpCodes(ctx context.Context, tx *sql.Tx) error {
	query := `
		DROP TABLE otp_codes;

		DROP TYPE identifier_type_enum;
		
		DROP TYPE otp_type_enum;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
