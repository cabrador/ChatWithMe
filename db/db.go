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
	SELECT content, author 
	FROM messages 
	    INNER JOIN authors on authors.id = messages.author_id 
	WHERE messages.user_id=$1 AND messages.persona_id=$2 ORDER BY messages.order_number
`

func (db *Database) GetUserPersonaMessages(userId, personaId int) ([]Message, error) {
	var msg []Message
	err := db.db.Select(&msg, getMessages, userId, personaId)
	return msg, err
}
