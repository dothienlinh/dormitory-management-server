package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateRoomRulesTable, downCreateRoomRulesTable)
}

func upCreateRoomRulesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS room_rules (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		category VARCHAR(50) NOT NULL CHECK (category IN ('general', 'safety', 'hygiene', 'behavior')),
		priority INTEGER DEFAULT 1,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL
	);

	CREATE INDEX IF NOT EXISTS idx_room_rules_category ON room_rules(category);
	CREATE INDEX IF NOT EXISTS idx_room_rules_is_active ON room_rules(is_active);
	CREATE INDEX IF NOT EXISTS idx_room_rules_priority ON room_rules(priority);
	CREATE INDEX IF NOT EXISTS idx_room_rules_deleted_at ON room_rules(deleted_at);
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateRoomRulesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS room_rules;")
	return err
}
