package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateCleaningSchedulesTable, downCreateCleaningSchedulesTable)
}

func upCreateCleaningSchedulesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS cleaning_schedules (
		id SERIAL PRIMARY KEY,
		room_id INTEGER NULL,
		day_of_week VARCHAR(20) NOT NULL CHECK (day_of_week IN ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday')),
		start_time VARCHAR(10) NOT NULL,
		end_time VARCHAR(10) NOT NULL,
		type VARCHAR(20) NOT NULL CHECK (type IN ('regular', 'deep', 'inspection')),
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL,
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_cleaning_schedules_room_id ON cleaning_schedules(room_id);
	CREATE INDEX IF NOT EXISTS idx_cleaning_schedules_day_of_week ON cleaning_schedules(day_of_week);
	CREATE INDEX IF NOT EXISTS idx_cleaning_schedules_type ON cleaning_schedules(type);
	CREATE INDEX IF NOT EXISTS idx_cleaning_schedules_is_active ON cleaning_schedules(is_active);
	CREATE INDEX IF NOT EXISTS idx_cleaning_schedules_deleted_at ON cleaning_schedules(deleted_at);
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateCleaningSchedulesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS cleaning_schedules;")
	return err
}
