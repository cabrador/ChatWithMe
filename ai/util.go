package ai

import "chatwithme/db"

type openAiRole string

const (
	userRole      openAiRole = "user"
	assistantRole openAiRole = "assistant"
)

type chatRequest struct {
	Model    string       `json:"model"`
	User     string       `json:"user"`
	Messages []db.Message `json:"messages"`
}
