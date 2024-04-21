package types

import "fmt"

type Message struct {
	UserId      int    `db:"user_id"`
	PersonaId   int    `db:"persona_id"`
	Author      string `json:"role"`
	AuthorId    int    `json:"-" db:"author_id"`
	Content     string `json:"content" db:"content"`
	OrderNumber int    `db:"order_number"`
}

type Messages []Message

func (m Messages) String() string {
	var str string
	for i := 0; i < len(m); i = i + 2 {
		str = fmt.Sprintf("%s\n User: %s\n Clara: %s\n", str, m[i].Content, m[i+1].Content)
	}

	return str
}
