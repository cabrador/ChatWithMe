package ai

import (
	"chatwithme/types"
)

type openAiRole string

const (
	userRole      openAiRole = "user"
	assistantRole openAiRole = "assistant"
)

type chatRequest struct {
	Model    string          `json:"model"`
	User     string          `json:"user"`
	Messages []types.Message `json:"messages"`
}
