package main

import (
	"Schwarz--Internship--2026/services/reservation-base/main/proto" // Asigură-te că path-ul e corect
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errNoRowsDeleted = errors.New("no rows deleted")

func InsertReservation(ctx context.Context, db *sql.DB, res *proto.Reservation) (int64, error) {
	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("reservations").
		Columns("property_id", "user_id", "owner_id", "check_in_date", "check_out_date", "status").
		Values(
			res.PropertyId,
			res.UserId,
			res.OwnerId,
			res.CheckInDate,
			res.CheckOutDate,
			res.Status.String(), 
		).
		Suffix("RETURNING id")

	var id int64
	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectReservation(ctx context.Context, db *sql.DB, id int64) (*proto.Reservation, error) {
	var res proto.Reservation
	var checkIn, checkOut, createdAt, updatedAt time.Time
	var statusStr string

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select(
		"id",
		"property_id",
		"user_id",
		"owner_id",
		"check_in_date",
		"check_out_date",
		"status",
		"created_at",
		"updated_at").
		From("reservations").
		Where(sq.Eq{"id": id})

	err := query.RunWith(db).QueryRowContext(ctx).Scan(
		&res.Id,
		&res.PropertyId,
		&res.UserId,
		&res.OwnerId,
		&checkIn,
		&checkOut,
		&statusStr,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	res.CheckInDate = checkIn.Format("2006-01-02")
	res.CheckOutDate = checkOut.Format("2006-01-02")
	res.Status = proto.ReservationStatus(proto.ReservationStatus_value[statusStr])
	res.CreatedAt = timestamppb.New(createdAt)
	res.UpdatedAt = timestamppb.New(updatedAt)

	return &res, nil
}

func SelectReservationListInDB(ctx context.Context, db *sql.DB, offsetID int64, pageSize int64, filterPropertyIds []int64, filterUserId int64, filterOwnerId int64, filterStatus proto.ReservationStatus) ([]*proto.Reservation, error) {
	base := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select(
		"id",
		"property_id",
		"user_id",
		"owner_id",
		"check_in_date",
		"check_out_date",
		"status",
		"created_at",
		"updated_at").
		From("reservations")

	if len(filterPropertyIds) > 0 {
		base = base.Where(sq.Eq{"property_id": filterPropertyIds})
	}
	if filterUserId > 0 {
		base = base.Where(sq.Eq{"user_id": filterUserId})
	}
	if filterOwnerId > 0 {
		base = base.Where(sq.Eq{"owner_id": filterOwnerId})
	}
	if filterStatus != proto.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED {
		base = base.Where(sq.Eq{"status": filterStatus.String()})
	}

	if offsetID > 0 {
		base = base.Where(sq.Gt{"id": offsetID})
	}

	base = base.OrderBy("id ASC").Suffix("LIMIT ?", uint64(pageSize+1))

	query, args, _ := base.ToSql()
	slog.Info("Executing list query", "query", query, "args", args)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var reservations []*proto.Reservation
	for rows.Next() {
		var r proto.Reservation
		var checkIn, checkOut, createdAt, updatedAt time.Time
		var statusStr string

		if err := rows.Scan(&r.Id, &r.PropertyId, &r.UserId, &r.OwnerId, &checkIn, &checkOut, &statusStr, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan reservation: %w", err)
		}

		r.CheckInDate = checkIn.Format("2006-01-02")
		r.CheckOutDate = checkOut.Format("2006-01-02")
		r.Status = proto.ReservationStatus(proto.ReservationStatus_value[statusStr])
		r.CreatedAt = timestamppb.New(createdAt)
		r.UpdatedAt = timestamppb.New(updatedAt)

		reservations = append(reservations, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading database rows: %w", err)
	}

	return reservations, nil
}

func DeleteReservationInDB(ctx context.Context, db *sql.DB, id int64) error {
	del := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("reservations").
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
		return errNoRowsDeleted
	}

	return nil
}