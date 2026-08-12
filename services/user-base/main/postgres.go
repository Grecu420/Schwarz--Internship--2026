package main

import (
	"Schwarz--Internship--2026/services/user-base/main/proto"
	"database/sql"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func SelectUser(db *sql.DB, id int64) (*proto.User, error) {

	var resUser proto.User
	resUser.Id = id

	query := `
		SELECT first_name, last_name, user_name, email, hashed_password, created_at
		FROM users
		WHERE id = $1
	`
	var created_at time.Time

	err := db.QueryRow(query, id).Scan(
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

func InsertUser(db *sql.DB, user *proto.User) (int64, error) {
	// execute query + get new id
	query := `
		INSERT INTO users (id, first_name, last_name, user_name, email, hashed_password, created_at)
		VALUES (DEFAULT, $1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int64
	err := db.QueryRow(query,
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
