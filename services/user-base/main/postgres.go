package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserFilter struct {
	FirstName *string
	LastName  *string
}

func SelectUser(ctx context.Context, db *sql.DB, email string) (*proto.User, error) {
	var resUser proto.User
	var created_at time.Time

	user := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select(
		"id",
		"first_name",
		"last_name",
		"user_name",
		"email",
		"hashed_password",
		"created_at").
		From("users").
		Where(sq.Eq{"email": email})

	err := user.RunWith(db).QueryRowContext(ctx).Scan(
		&resUser.Id,
		&resUser.FirstName,
		&resUser.LastName,
		&resUser.UserName,
		&resUser.Email,
		&resUser.Password,
		&created_at)

	if err != nil {
		return nil, err
	}

	resUser.CreatedAt = timestamppb.New(created_at)

	return &resUser, nil
}

func InsertUser(ctx context.Context, db *sql.DB, user *proto.User) (int64, error) {
	// execute query + get new id

	ins := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("users").
		Columns("id", "first_name", "last_name", "user_name", "email", "hashed_password", "created_at").
		Values(sq.Expr("DEFAULT"), user.FirstName, user.LastName, user.UserName, user.Email, user.Password, user.CreatedAt.AsTime()).
		Suffix("RETURNING id")
	// query := `
	// 	INSERT INTO users (id, first_name, last_name, user_name, email, hashed_password, created_at)
	// 	VALUES (DEFAULT, $1, $2, $3, $4, $5, $6)
	// 	RETURNING id
	// `
	var id int64
	// err := db.QueryRowContext(ctx,
	// 	query,
	// 	user.FirstName,
	// 	user.LastName,
	// 	user.UserName,
	// 	user.Email,
	// 	user.Password,
	// 	user.CreatedAt.AsTime()).Scan(&id)

	err := ins.RunWith(db).QueryRowContext(ctx).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectUserListInDB(ctx context.Context, db *sql.DB, offsetID int64, pageSize int64, firstName string, lastName string) ([]*proto.User, error) {

	base := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("id", "first_name", "last_name", "user_name", "email").From("users")

	if firstName != "" {
		base = base.Where(sq.Eq{"first_name": firstName})
	}

	if lastName != "" {
		base = base.Where(sq.Eq{"last_name": lastName})
	}

	if offsetID > 0 {
		base = base.Where(sq.Gt{"id": offsetID})
	}

	base = base.OrderBy("id ASC").Suffix("LIMIT ?", uint64(pageSize+1))

	query, args, _ := base.ToSql()
	slog.Info("base", "query", query)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}
	defer rows.Close()

	var users []*proto.User
	for rows.Next() {
		var u proto.User
		if err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.UserName, &u.Email); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading database rows: %w", err)
	}

	return users, nil
}
