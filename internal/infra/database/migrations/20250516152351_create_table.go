package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTable, downCreateTable)
}

func upCreateTable(ctx context.Context, tx *sql.Tx) error {
	createTypeQuery := `
	CREATE TYPE user_role AS ENUM ('student', 'admin', 'staff');
	CREATE TYPE user_gender AS ENUM ('male', 'female', 'other');
	CREATE TYPE user_status AS ENUM ('active', 'inactive');
	CREATE TYPE room_status AS ENUM ('available', 'occupied', 'maintenance');
	CREATE TYPE contract_status AS ENUM ('pending', 'active', 'expired', 'cancelled');
	`

	if _, err := tx.ExecContext(ctx, createTypeQuery); err != nil {
		return err
	}

	createTableAndIndexQuery := `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255),
		student_code VARCHAR(255) UNIQUE,
		email VARCHAR(255) UNIQUE,
		password TEXT,
		role user_role DEFAULT 'student',
		gender user_gender DEFAULT 'other',
		status user_status DEFAULT 'active',
		phone VARCHAR(255),
		birthday TIMESTAMPTZ,
		avatar TEXT,
		room_rent_id INTEGER,
		contract_id INTEGER,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);
	CREATE INDEX idx_users_email ON users (email);
	CREATE INDEX idx_users_student_code ON users (student_code);

	CREATE TABLE room_categories (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		capacity INTEGER NOT NULL CHECK (capacity IN (8, 6, 4)),
		price INTEGER NOT NULL,
		acreage INTEGER NOT NULL,
		description TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

	CREATE TABLE rooms (
		id SERIAL PRIMARY KEY,
		room_number TEXT NOT NULL UNIQUE,
		status room_status NOT NULL,
		room_category_id INTEGER NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);
	CREATE INDEX idx_rooms_room_number ON rooms (room_number);

	CREATE TABLE contracts (
		id SERIAL PRIMARY KEY,
		room_id INTEGER NOT NULL,
		user_id INTEGER,
		start_date TIMESTAMPTZ NOT NULL,
		end_date TIMESTAMPTZ NOT NULL,
		status contract_status NOT NULL,
		code VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);
	CREATE INDEX idx_contracts_code ON contracts (code);

	CREATE TABLE room_rents (
		id SERIAL PRIMARY KEY,
		room_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);
	`

	if _, err := tx.ExecContext(ctx, createTableAndIndexQuery); err != nil {
		return err
	}

	addForeignKeyConstraintsQuery := `
	ALTER TABLE users 
		ADD CONSTRAINT fk_contracts_user FOREIGN KEY (contract_id) REFERENCES contracts(id),
		ADD CONSTRAINT fk_room_rents_user FOREIGN KEY (room_rent_id) REFERENCES room_rents(id);

	ALTER TABLE rooms
		ADD CONSTRAINT fk_room_categories_rooms FOREIGN KEY (room_category_id) REFERENCES room_categories(id);

	ALTER TABLE contracts
		ADD CONSTRAINT fk_users_contract FOREIGN KEY (user_id) REFERENCES users(id),
		ADD CONSTRAINT fk_rooms_contract FOREIGN KEY (room_id) REFERENCES rooms(id);

	ALTER TABLE room_rents
		ADD CONSTRAINT fk_users_room_rents FOREIGN KEY (user_id) REFERENCES users(id),
		ADD CONSTRAINT fk_rooms_room_rents FOREIGN KEY (room_id) REFERENCES rooms(id);
	`

	if _, err := tx.ExecContext(ctx, addForeignKeyConstraintsQuery); err != nil {
		return err
	}

	return nil
}

func downCreateTable(ctx context.Context, tx *sql.Tx) error {
	dropForeignKeyConstraintsQuery := `
	ALTER TABLE users 
		DROP CONSTRAINT IF EXISTS fk_contracts_user,
		DROP CONSTRAINT IF EXISTS fk_room_rents_user;

	ALTER TABLE rooms
		DROP CONSTRAINT IF EXISTS fk_room_categories_rooms;

	ALTER TABLE contracts
		DROP CONSTRAINT IF EXISTS fk_users_contract,
		DROP CONSTRAINT IF EXISTS fk_rooms_contract;

	ALTER TABLE room_rents
		DROP CONSTRAINT IF EXISTS fk_users_room_rents,
		DROP CONSTRAINT IF EXISTS fk_rooms_room_rents;
	`

	if _, err := tx.ExecContext(ctx, dropForeignKeyConstraintsQuery); err != nil {
		return err
	}

	dropTables := `
	DROP TABLE IF EXISTS room_rents;
	DROP TABLE IF EXISTS contracts;
	DROP TABLE IF EXISTS rooms;
	DROP TABLE IF EXISTS room_categories;
	DROP TABLE IF EXISTS users;
	`

	if _, err := tx.ExecContext(ctx, dropTables); err != nil {
		return err
	}

	dropTypesQuery := `
	DROP TYPE IF EXISTS contract_status;
	DROP TYPE IF EXISTS room_status;
	DROP TYPE IF EXISTS user_status;
	DROP TYPE IF EXISTS user_gender;
	DROP TYPE IF EXISTS user_role;
	`

	if _, err := tx.ExecContext(ctx, dropTypesQuery); err != nil {
		return err
	}

	return nil
}
