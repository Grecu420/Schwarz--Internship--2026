package main

import (
	"Schwarz--Internship--2026/services/api-rest-gateway/main/proto"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func InsertFriendRequestInDB(ctx context.Context, db *sql.DB, id string, senderID string, receiverID string, status int32) error {

	query := `
		INSERT INTO friend_requests (id, sender_id, receiver_id, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.ExecContext(ctx, query, id, senderID, receiverID, status)

	return err
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
