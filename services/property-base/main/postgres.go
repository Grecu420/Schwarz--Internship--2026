package main

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func InsertPropertyInDB(ctx context.Context, db *sql.DB, name string, description string, user_id int64, address string, price int, lng, lat float64) (int64, error) {
	var generatedID int64

	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("properties").
		Columns("user_id", "name", "description", "address", "location", "price").
		Values(user_id, name, description, address, sq.Expr("ST_MakePoint(?, ?)::geography", lng, lat), price)

	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&generatedID)

	return generatedID, err
}
