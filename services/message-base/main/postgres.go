package main

import (
	"Schwarz--Internship--2026/services/message-base/main/proto"
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func InsertMessage(ctx context.Context, db *sql.DB, conversationID int64, senderID int64, content string) (int64, time.Time, error) {
	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("messages").
		Columns("conversation_id", "sender_id", "content").
		Values(conversationID, senderID, content).
		Suffix("RETURNING id, created_at")

	var id int64
	var createdAt time.Time

	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&id, &createdAt)
	if err != nil {
		return 0, time.Time{}, err
	}

	return id, createdAt, nil
}

func SelectMessages(ctx context.Context, db *sql.DB, conversationID int64) ([]*proto.Message, error) {

	sel := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "conversation_id", "sender_id", "content", "created_at").
		From("messages").
		Where(sq.Eq{"conversation_id": conversationID}).
		OrderBy("created_at DESC")
	rows, err := sel.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}
	defer rows.Close()

	var messages []*proto.Message
	for rows.Next() {
		// Read each row
		var message proto.Message
		var createdAt time.Time

		if err := rows.Scan(
			&message.Id,
			&message.ConversationId,
			&message.SenderId,
			&message.Content,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("scan friend request: %w", err)
		}
		message.CreatedAt = timestamppb.New(createdAt)
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading database rows: %w", err)
	}

	return messages, nil
}
