package database

import (
	"Schwarz--Internship--2026/services/common"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	bin, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, fmt.Errorf("failed to read db password secret: %w", err)
	}
	password := strings.TrimSpace(string(bin))

	user, err := common.GetRequiredEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}
	port, err := common.GetRequiredEnv("POSTGRES_PORT")
	if err != nil {
		return nil, err
	}
	dbName, err := common.GetRequiredEnv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}
	service, err := common.GetRequiredEnv("POSTGRES_SERVICE")
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(service, port),
		Path:   dbName,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()

	db, err := sql.Open("postgres", u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db driver: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close() // Prevent resource leak on failed connection
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
