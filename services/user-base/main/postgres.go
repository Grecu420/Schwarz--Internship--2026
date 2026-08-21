package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserFilter struct {
	FirstName *string
	LastName  *string
}

func SelectUser(ctx context.Context, db *sql.DB, email string) (*proto.User, error) {
	var resUser proto.User
	var created_at time.Time

	query := `
		SELECT id, first_name, last_name, user_name, email, hashed_password, created_at
		FROM users
		WHERE email = $1
	`

	err := db.QueryRowContext(ctx, query, email).Scan(
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
	query := `
		INSERT INTO users (id, first_name, last_name, user_name, email, hashed_password, created_at)
		VALUES (DEFAULT, $1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int64
	err := db.QueryRowContext(ctx,
		query,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Email,
		user.Password,
		user.CreatedAt.AsTime()).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectUserListInDB(ctx context.Context, db *sql.DB, offsetID int64, pageSize int64, firstName string, lastName string) ([]*proto.User, error) {

	baseQuery := `SELECT id, first_name, last_name, user_name, email FROM users`

	var whereClauses []string
	var args []any
	argNum := 1

	if firstName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("first_name = $%d", argNum))
		args = append(args, firstName)
		argNum++
	}

	if lastName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("last_name = $%d", argNum))
		args = append(args, lastName)
		argNum++
	}

	if offsetID > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("id > $%d", argNum))
		args = append(args, offsetID)
		argNum++
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY id ASC LIMIT $%d;", argNum)
	args = append(args, pageSize+1)

	rows, err := db.QueryContext(ctx, baseQuery, args...)
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
