package main

import (
	"Schwarz--Internship--2026/services/property-base/main/proto"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func InsertPropertyInDB(ctx context.Context, db *sql.DB, name string, description string, user_id int64, address string, price int, lng, lat float64, image_urls []string) (int64, error) {
	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("properties").
		Columns("user_id", "name", "description", "address", "location", "price", "image_urls").
		Values(user_id, name, description, address, sq.Expr("ST_MakePoint(?, ?)::geography", lng, lat), price, pq.Array(image_urls)).
		Suffix("RETURNING id")
	var generatedID int64
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
			"ST_Y(location::geometry) AS lat",
			"image_urls",
		).
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
		pq.Array(&property.ImageUrls),
	)

	if err != nil {
		return nil, err
	}

	return property, nil
}

func SelectPropertyListInDB(ctx context.Context, db *sql.DB, offsetID int64, page_size int64,
	owner *proto.FilterByOwner,
	name *proto.FilterByName,
	price *proto.FilterByPriceRange,
	location *proto.FilterByLocation) ([]*proto.Property, error) {
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
			"ST_Y(location::geometry) AS lat",
			"image_urls").
		From("properties")

	// Add filter for offset (starting ID)
	base = base.Where(sq.GtOrEq{"id": offsetID})
	if owner != nil {
		base = base.Where(sq.Eq{"user_id": owner.Value})
	}
	if name != nil {
		base = base.Where(sq.Like{"name": "%" + name.Value + "%"})
	}
	if price != nil {
		if price.GetMin() > 0 {
			base = base.Where(sq.GtOrEq{"price": price.Min})
		}
		if price.GetMax() > 0 {
			base = base.Where(sq.LtOrEq{"price": price.Max})
		}
	}

	if location != nil && location.Center != nil {
		base = base.Where(sq.Expr("ST_DWithin(location, ST_MakePoint(?, ?)::geography, ?)",
			location.Center.Long, location.Center.Lat, location.Radius))
	}

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
			pq.Array(&property.ImageUrls),
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

func CountPropertiesInDB(ctx context.Context, db *sql.DB,
	owner *proto.FilterByOwner,
	name *proto.FilterByName,
	price *proto.FilterByPriceRange,
	location *proto.FilterByLocation) (int64, error) {
	// Build query
	base := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From("properties")

	if owner != nil {
		base = base.Where(sq.Eq{"user_id": owner.Value})
	}
	if name != nil {
		base = base.Where(sq.Like{"name": "%" + name.Value + "%"})
	}
	if price != nil {
		if price.GetMin() > 0 {
			base = base.Where(sq.GtOrEq{"price": price.Min})
		}
		if price.GetMax() > 0 {
			base = base.Where(sq.LtOrEq{"price": price.Max})
		}
	}

	if location != nil && location.Center != nil {
		base = base.Where(sq.Expr("ST_DWithin(location, ST_MakePoint(?, ?)::geography, ?)",
			location.Center.Long, location.Center.Lat, location.Radius))
	}

	// Execute query
	var count int64
	err := base.RunWith(db).QueryRowContext(ctx).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
