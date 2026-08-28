package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
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
