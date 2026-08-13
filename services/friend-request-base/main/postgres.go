package main

import (
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func SelectFriendListInDB(ctx context.Context, db *sql.DB, offsetID int64, page_size int64, senderID string, receiverID string, status proto.RequestStatus) {

	// build query
	baseQuery := "SELECT id, name, status, role, created_at FROM users"
	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var whereClauses []string
	var args []any
	argNum := 1

	// Optional Filter 1: Status
	if req.GetStatusFilter() != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argNum))
		args = append(args, req.GetStatusFilter())
		argNum++
	}

	// Optional Filter 2: Role
	if req.GetRoleFilter() != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("role = $%d", argNum))
		args = append(args, req.GetRoleFilter())
		argNum++
	}

	// Optional Filter 3: Search Query (ILIKE / partial match)
	if req.GetSearchQuery() != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argNum))
		args = append(args, "%"+req.GetSearchQuery()+"%")
		argNum++
	}
}
