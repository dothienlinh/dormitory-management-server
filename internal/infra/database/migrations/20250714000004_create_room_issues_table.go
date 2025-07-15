package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateRoomIssuesTable, downCreateRoomIssuesTable)
}

func upCreateRoomIssuesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS room_issues (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		category VARCHAR(100) NOT NULL,
		priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
		status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'resolved', 'rejected')),
		room_id INTEGER NOT NULL,
		reported_by INTEGER NOT NULL,
		assigned_to INTEGER NULL,
		resolved_date TIMESTAMP NULL,
		resolution TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL,
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
		FOREIGN KEY (reported_by) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE INDEX IF NOT EXISTS idx_room_issues_room_id ON room_issues(room_id);
	CREATE INDEX IF NOT EXISTS idx_room_issues_reported_by ON room_issues(reported_by);
	CREATE INDEX IF NOT EXISTS idx_room_issues_status ON room_issues(status);
	CREATE INDEX IF NOT EXISTS idx_room_issues_category ON room_issues(category);
	CREATE INDEX IF NOT EXISTS idx_room_issues_created_at ON room_issues(created_at);
	CREATE INDEX IF NOT EXISTS idx_room_issues_deleted_at ON room_issues(deleted_at);
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateRoomIssuesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS room_issues;")
	return err
}
