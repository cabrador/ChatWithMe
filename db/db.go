package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/petr-hanzl/chatwithme/types"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	ssl               = "disabled"
	UserAuthorId      = 1
	AssistantAuthorId = 2
)

// test

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
	INSERT INTO messages(user_id, persona_id, content, order_number, author_id)
	VALUES (:user_id, :persona_id, :content, :order_number, :author_id);
`

func (db *Database) InsertMessages(msgs []types.Message) (sql.Result, error) {
	return db.db.NamedExec(insertMessage, msgs)
}

const getMessages = `
	SELECT user_id, persona_id, content, author, order_number
	FROM messages 
	    INNER JOIN authors on authors.id = messages.author_id 
	WHERE messages.user_id=$1 AND messages.persona_id=$2 ORDER BY messages.order_number
`

func (db *Database) GetUserPersonaMessages(userId, personaId int) ([]types.Message, error) {
	var msg []types.Message
	err := db.db.Select(&msg, getMessages, userId, personaId)
	return msg, err
}
