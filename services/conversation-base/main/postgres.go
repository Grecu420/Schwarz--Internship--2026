package main

import (
	"Schwarz--Internship--2026/services/conversation-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrConversationExists = errors.New("conversation already exists")

func InsertConversation(ctx context.Context, db *sql.DB, user1ID int64, user2ID int64) (int64, error) {
	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("conversations").
		Columns("user1_id", "user2_id").
		Values(user1ID, user2ID).
		Suffix("RETURNING id")

	var id int64

	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&id)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			slog.Info("conversation already exists")
			return 0, ErrConversationExists
		}
		return 0, err
	}

	return id, nil
}

func SelectConversationsListInDB(ctx context.Context, db *sql.DB, userID int64) ([]*proto.Conversation, error) {
	sel := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "user1_id", "user2_id", "created_at", "updated_at").
		From("conversations").
		Where(sq.Or{
			sq.Eq{"user1_id": userID},
			sq.Eq{"user2_id": userID},
		}).
		OrderBy("updated_at DESC")

	rows, err := sel.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}
	defer rows.Close()

	var conversations []*proto.Conversation
	for rows.Next() {
		var conv proto.Conversation
		var createdAt time.Time
		var updatedAt time.Time

		if err := rows.Scan(
			&conv.Id,
			&conv.User1Id,
			&conv.User2Id,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan conversation: %w", err)
		}
		
		conv.CreatedAt = timestamppb.New(createdAt)
		conv.UpdatedAt = timestamppb.New(updatedAt)
		conversations = append(conversations, &conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading database rows: %w", err)
	}

	return conversations, nil
}