package main

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

func InsertPropertyInDB(ctx context.Context, db *sql.DB, name string, description string, user_id int64, address string, price int, lng, lat float64) (int64, error) {
	var generatedID int64

	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("properties").
		Columns("user_id", "name", "description", "address", "location", "price").
		Values(user_id, name, description, address, sq.Expr("ST_MakePoint(?, ?)::geography", lng, lat), price).
		Suffix("RETURNING id")
	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&generatedID)

	return generatedID, err
}

var NoRowsDeleted = errors.New("no rows deleted")

func DeletePropertyInDB(ctx context.Context, db *sql.DB, id int64) error {

	del := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("properties").
		Where(sq.Eq{"id": id})

	res, err := del.RunWith(db).ExecContext(ctx)

	if err != nil {
		return err
	}

	r, err := res.RowsAffected()

	if err != nil {
		return err
	}
	if r == 0 {
		return NoRowsDeleted
	}
	return nil

}
