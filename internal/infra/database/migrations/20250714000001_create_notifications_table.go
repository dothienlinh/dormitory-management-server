package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateNotificationsTable, downCreateNotificationsTable)
}

func upCreateNotificationsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT,
		type VARCHAR(50) NOT NULL CHECK (type IN ('payment', 'maintenance', 'event', 'service', 'announcement')),
		priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
		is_read BOOLEAN DEFAULT FALSE,
		user_id INTEGER NOT NULL,
		read_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
	CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
	CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
	CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications(deleted_at);
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateNotificationsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS notifications;")
	return err
}
