package main

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func InsertConversation(ctx context.Context, db *sql.DB, user1ID int64, user2ID int64) (int64, error) {
	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("conversations").
		Columns("user1_id", "user2_id").
		Values(user1ID, user2ID).
		Suffix("RETURNING id")

	var id int64

	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
