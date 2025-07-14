package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateEventsTable, downCreateEventsTable)
}

func upCreateEventsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		event_date DATE NOT NULL,
		start_time VARCHAR(10) NOT NULL,
		end_time VARCHAR(10) NOT NULL,
		location VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL CHECK (type IN ('deadline', 'maintenance', 'event', 'meeting')),
		is_mandatory BOOLEAN DEFAULT FALSE,
		organizer VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL
	);

	CREATE INDEX IF NOT EXISTS idx_events_event_date ON events(event_date);
	CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
	CREATE INDEX IF NOT EXISTS idx_events_is_mandatory ON events(is_mandatory);
	CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);
	CREATE INDEX IF NOT EXISTS idx_events_deleted_at ON events(deleted_at);
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateEventsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS events;")
	return err
}
