package main

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func InsertMessage(ctx context.Context, db *sql.DB, conversationID int64, senderID int64, content string) (int64, time.Time, error) {
	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("messages").
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