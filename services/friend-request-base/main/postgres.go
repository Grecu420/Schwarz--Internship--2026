package main

import (
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
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

func SelectFriendListInDB(ctx context.Context, db *sql.DB, offsetID int64, page_size int64, senderID string, receiverID string, status proto.RequestStatus) ([]*proto.FriendRequest, error) {

	// Build query
	baseQuery := `SELECT "id", "sender_id", "receiver_id", "status", "created_at" FROM "friend_requests"`

	var whereClauses []string
	var args []any
	argNum := 1

	// Optional Filter 1: senderID
	if senderID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("sender_id = $%d", argNum))
		args = append(args, senderID)
		argNum++
	}

	// Optional Filter 2: receiverID
	if receiverID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("receiver_id = $%d", argNum))
		args = append(args, senderID)
		argNum++
	}

	// Optional Filter 3: status
	if status != proto.RequestStatus_STATUS_UNKNOWN {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argNum))
		args = append(args, status)
		argNum++
	}

	// Add filter for offset (starting ID)
	whereClauses = append(whereClauses, fmt.Sprintf("id >= $%d", argNum))
	args = append(args, offsetID)
	argNum++

	// Append where clauses
	baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")

	// Append ordering and limit to page size
	baseQuery += fmt.Sprintf(" ORDER BY id LIMIT $%d;", argNum)
	args = append(args, page_size+1)

	// Execute query
	rows, err := db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}
	defer rows.Close()

	// Scan rows
	var friendRequests []*proto.FriendRequest
	for rows.Next() {
		var fr proto.FriendRequest
		var createdAt time.Time

		if err := rows.Scan(&fr.Id, &fr.SenderId, &fr.ReceiverId, &fr.Status, &createdAt); err != nil {
			return nil, fmt.Errorf("scan friend request: %w", err)
		}
		fr.CreatedAt = timestamppb.New(createdAt)
		friendRequests = append(friendRequests, &fr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading database rows: %w", err)
	}

	return friendRequests, nil
}
