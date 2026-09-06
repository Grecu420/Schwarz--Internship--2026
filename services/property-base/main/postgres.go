package main

import (
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func SelectPropertyInDB(ctx context.Context, db *sql.DB, id int64) (*proto.Property, error) {

	var property = &proto.Property{Location: &proto.Location{}}

	build := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"user_id",
			"name",
			"description",
			"address",
			"price",
			"ST_X(location::geometry) AS lng",
			"ST_Y(location::geometry) AS lat").
		From("properties").
		Where(sq.Eq{"id": id})

	err := build.RunWith(db).QueryRowContext(ctx).Scan(
		&property.Id,
		&property.UserId,
		&property.Name,
		&property.Description,
		&property.Address,
		&property.Price,
		&property.Location.Long,
		&property.Location.Lat,
	)

	if err != nil {
		return nil, err
	}

	return property, nil
}

func SelectPropertyListInDB(ctx context.Context, db *sql.DB, offsetID int64, page_size int64, userID int64) ([]*proto.Property, error) {
	// Build query
	base := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id",
			"user_id",
			"name",
			"description",
			"address",
			"price",
			"ST_X(location::geometry) AS lng",
			"ST_Y(location::geometry) AS lat").
		From("properties").
		Where(sq.Eq{"user_id": userID})

	// Add filter for offset (starting ID)
	base = base.Where(sq.GtOrEq{"id": offsetID})

	// Append ordering and limit
	base = base.OrderBy("id ASC").
		Limit(uint64(page_size + 1))

	// Execute query
	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}
	defer rows.Close()

	// Scan rows
	var properties []*proto.Property
	for rows.Next() {
		var property = &proto.Property{Location: &proto.Location{}}

		if err := rows.Scan(
			&property.Id,
			&property.UserId,
			&property.Name,
			&property.Description,
			&property.Address,
			&property.Price,
			&property.Location.Long,
			&property.Location.Lat,
		); err != nil {
			return nil, fmt.Errorf("scan property: %w", err)
		}

		properties = append(properties, property)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading database rows: %w", err)
	}

	return properties, nil
}
