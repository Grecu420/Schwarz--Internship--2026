package main

import (
	"Schwarz--Internship--2026/services/friend-request-base/main/proto"
	"context"
	"database/sql"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
)

func InsertFriendRequestInDB(ctx context.Context, db *sql.DB, senderID string, receiverID string, status int32) (int64, error) {
	var generatedID int64

	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("friend_requests").
		Columns("sender_id", "receiver_id", "status").
		Values(senderID, receiverID, status).
		Suffix("RETURNING id")

	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&generatedID)

	return generatedID, err
}

func UpdateFriendRequestStatusInDB(ctx context.Context, db *sql.DB, id int64, newStatus proto.RequestStatus) (int64, error) {
	upd := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("friend_requests").
		Set("status", newStatus).
		Where(sq.Eq{"id": id})

	res, err := upd.RunWith(db).ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func SelectFriendListInDB(ctx context.Context, db *sql.DB, offsetID int64, page_size int64, senderID string, receiverID string, status proto.RequestStatus) ([]*proto.FriendRequest, error) {

	// Build query

	base := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "sender_id", "receiver_id", "status", "created_at").
		From("friend_requests")

	// Optional Filter 1: senderID
	if senderID != "" {
		base = base.Where(sq.Eq{"sender_id": senderID})
	}

	// Optional Filter 2: receiverID
	if receiverID != "" {
		base = base.Where(sq.Eq{"receiver_id": receiverID})
	}

	// Optional Filter 3: status
	if status != proto.RequestStatus_STATUS_UNKNOWN {
		base = base.Where(sq.Eq{"status": status})
	}

	// Add filter for offset (starting ID)
	base = base.Where(sq.GtOrEq{"id": offsetID})

	// Append ordering and limit to page size
	base = base.OrderBy("id").Suffix("LIMIT ?", uint64(page_size+1))

	// Execute query
	rows, err := base.RunWith(db).QueryContext(ctx)
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
