package main

import (
	"Schwarz--Internship--2026/services/property-list-base/main/proto"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func applyPropertyFilters(
	base sq.SelectBuilder,
	owner *proto.FilterByOwner,
	name *proto.FilterByName,
	price *proto.FilterByPriceRange,
	location *proto.FilterByLocation,
	dateInterval *proto.FilterByDateInterval,
) sq.SelectBuilder {
	if owner != nil {
		base = base.Where(sq.Eq{"properties.user_id": owner.Value})
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
	if dateInterval != nil && dateInterval.GetCheckInDate() != "" && dateInterval.GetCheckOutDate() != "" {
		base = base.LeftJoin(
			"reservations ON reservations.property_id = properties.id AND reservations.check_in_date < ? AND reservations.check_out_date > ? AND reservations.status = 'RESERVATION_STATUS_CONFIRMED'",
			dateInterval.GetCheckOutDate(),
			dateInterval.GetCheckInDate(),
		)
		base = base.Where(sq.Eq{"reservations.id": nil})
	}
	return base
}

func SelectPropertyListInDB(ctx context.Context, db *sql.DB, offsetID int64, page_size int64,
	sort_type proto.SortType,
	owner *proto.FilterByOwner,
	name *proto.FilterByName,
	price *proto.FilterByPriceRange,
	location *proto.FilterByLocation,
	dateInterval *proto.FilterByDateInterval) ([]*proto.Property, error) {

	base := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"properties.id",
			"properties.user_id",
			"name",
			"description",
			"address",
			"price",
			"ST_X(location::geometry) AS lng",
			"ST_Y(location::geometry) AS lat",
			"image_urls").
		From("properties").
		Where(sq.GtOrEq{"properties.id": offsetID})

	base = applyPropertyFilters(base, owner, name, price, location, dateInterval)

	switch sort_type {
	case proto.SortType_SORT_TYPE_NAME:
		base = base.OrderBy("name ASC")
	case proto.SortType_SORT_TYPE_PRICE_ASC:
		base = base.OrderBy("price ASC")
	case proto.SortType_SORT_TYPE_PRICE_DESC:
		base = base.OrderBy("price DESC")
	default:
		base = base.OrderBy("properties.id ASC")
	}
	base = base.Limit(uint64(page_size + 1))

	rows, err := base.RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}
	defer rows.Close()

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
	location *proto.FilterByLocation,
	dateInterval *proto.FilterByDateInterval) (int64, error) {

	base := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("COUNT(*)").
		From("properties")

	base = applyPropertyFilters(base, owner, name, price, location, dateInterval)

	var count int64
	err := base.RunWith(db).QueryRowContext(ctx).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
