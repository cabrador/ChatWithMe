package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const ssl = "disabled"

func MakeDb() (*Database, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	db.MustExec(schema)

	return &Database{db: db}, nil

}

type Database struct {
	db *sqlx.DB
}

const insertMessage = `
	INSERT INTO messages(user_id, persona_id, content, order_number)
	VALUES ($1, $2, $3, $4);
`

func (db *Database) InsertMessage(userId, personaId int, content string, orderNumber int) (sql.Result, error) {
	return db.db.Exec(insertMessage, userId, personaId, content, orderNumber)
}

const getMessages = `
	SELECT * FROM messages WHERE user_id=$1 AND persona_id=$2
`

func (db *Database) GetUserPersonaMessages(userId, personaId int) (sql.Result, error) {
	return db.db.Exec(getMessages, userId, personaId)
}
