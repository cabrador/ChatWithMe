package ai

import (
	"net/http"
)

func MakeChatGenerator() {}

type ChatGenerator interface {
	Generate(chatMessage) error
}

type chatGenerator struct {
	apiKey  string
	client  *http.Client
	history []chatMessage
}

func (c *chatGenerator) Generate(message chatMessage) error {
	c.history = append(c.history, message)

	//req := chatRequest{
	//	Model:    "gpt-3.5-turbo",
	//	User:     "idk",
	//	Messages: c.history,
	//}
	//req, err := http.NewRequest(http.MethodPost)
	//c.client.Do()

	return nil
}
