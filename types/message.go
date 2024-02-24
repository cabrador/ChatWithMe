package types

type Message struct {
	UserId      int    `db:"user_id"`
	PersonaId   int    `db:"persona_id"`
	Author      string `json:"role"`
	AuthorId    int    `json:"-" db:"author_id"`
	Content     string `json:"content" db:"content"`
	OrderNumber int    `db:"order_number"`
}
