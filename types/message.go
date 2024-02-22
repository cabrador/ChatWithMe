package types

type Message struct {
	Author  string `json:"role"`
	Content string `json:"content"`
}
