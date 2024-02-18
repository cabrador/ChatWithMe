package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const ssl = "disabled"

func Init() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	db.MustExec(schema)

	return db, nil

}
