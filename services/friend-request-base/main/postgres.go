package main

import (
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
	"context"
	"database/sql"
)

func InsertFriendRequestInDB(ctx context.Context, db *sql.DB, senderID string, receiverID string, status int32) (int64, error) {
	var generatedID int64

	query := `
		INSERT INTO friend_requests (sender_id, receiver_id, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := db.QueryRowContext(ctx, query, senderID, receiverID, status).Scan(&generatedID)

	return generatedID, err
}

func UpdateFriendRequestStatusInDB(ctx context.Context, db *sql.DB, id int64, newStatus proto.RequestStatus) (int64, error) {
	query := `
		UPDATE friend_requests
		SET status = $1
		WHERE id = $2
	`

	res, err := db.ExecContext(ctx, query, newStatus, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
