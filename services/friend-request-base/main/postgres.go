package main

import (
	"context"
	"database/sql"
)

func InsertFriendRequestInDB(ctx context.Context, db *sql.DB, id string, senderID string, receiverID string, status int32) error {

	query := `
		INSERT INTO friend_requests (id, sender_id, receiver_id, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.ExecContext(ctx, query, id, senderID, receiverID, status)

	return err
}
