package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateRoomAmenitiesTable, downCreateRoomAmenitiesTable)
}

func upCreateRoomAmenitiesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	CREATE TABLE room_amenities (
		id SERIAL PRIMARY KEY,
		room_id BIGINT NOT NULL,
		amenity_id BIGINT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ,
		UNIQUE (room_id, amenity_id),
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
		FOREIGN KEY (amenity_id) REFERENCES amenities(id) ON DELETE CASCADE
	);

	CREATE INDEX idx_room_amenities_room_id ON room_amenities(room_id);
	CREATE INDEX idx_room_amenities_amenity_id ON room_amenities(amenity_id);
	`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downCreateRoomAmenitiesTable(ctx context.Context, tx *sql.Tx) error {
	query := `
	DROP TABLE room_amenities;
	`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
