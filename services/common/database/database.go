package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func Connect() (*sql.DB, error) {
	bin, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, err
	}

	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, fmt.Errorf("missing port")
	}

	db, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return nil, fmt.Errorf("missing db")
	}
	service, ok := os.LookupEnv("POSTGRES_SERVICE")
	if !ok {
		return nil, fmt.Errorf("missing service")
	}

	return sql.Open("postgres", fmt.Sprintf("postgres://postgres:%s@%s:%s/%s?sslmode=disable", string(bin), service, port, db))
}

func Prepare() error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS users"); err != nil {
		return err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL, title VARCHAR)"); err != nil {
		return err
	}

	return nil
}
