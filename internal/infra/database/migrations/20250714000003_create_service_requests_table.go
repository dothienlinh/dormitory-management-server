package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateServiceRequestsTable, downCreateServiceRequestsTable)
}

func upCreateServiceRequestsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE IF NOT EXISTS service_requests (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		category VARCHAR(100) NOT NULL,
		priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
		status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'in_progress', 'completed', 'rejected')),
		user_id INTEGER NOT NULL,
		assigned_to INTEGER NULL,
		completion_date TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE INDEX IF NOT EXISTS idx_service_requests_user_id ON service_requests(user_id);
	CREATE INDEX IF NOT EXISTS idx_service_requests_assigned_to ON service_requests(assigned_to);
	CREATE INDEX IF NOT EXISTS idx_service_requests_status ON service_requests(status);
	CREATE INDEX IF NOT EXISTS idx_service_requests_category ON service_requests(category);
	CREATE INDEX IF NOT EXISTS idx_service_requests_priority ON service_requests(priority);
	CREATE INDEX IF NOT EXISTS idx_service_requests_created_at ON service_requests(created_at);
	CREATE INDEX IF NOT EXISTS idx_service_requests_deleted_at ON service_requests(deleted_at);
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateServiceRequestsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS service_requests;")
	return err
}
