package ai

import (
	"bytes"
	"chatwithme/db"
	"chatwithme/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	openAiApiUrl = "https://api.openai.com/v1/chat/completions"
)

type openAiResponse struct {
	Id                string
	Object            string
	Created           uint64
	Model             string
	Choices           []openAiChoices
	Usage             openAiUsage
	SystemFingerprint string `json:"system_fingerprint"`
}

type openAiChoices struct {
	Index        int
	Message      types.Message
	FinishReason string `json:"finish_reason"`
}

type openAiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func MakeChatGenerator(db *db.Database) ChatGenerator {
	return &chatGenerator{
		client: http.DefaultClient,
		apiKey: os.Getenv("OPENAI_API_TOKEN"),
		db:     db,
	}
}

type ChatGenerator interface {
	Generate(userId, personaId int, content string) ([]types.Message, error)
}

type chatGenerator struct {
	client *http.Client
	apiKey string
	db     *db.Database
}

func (c *chatGenerator) Generate(userId, personaId int, content string) ([]types.Message, error) {
	msgs, err := c.db.GetUserPersonaMessages(userId, personaId)
	if err != nil {
		return nil, fmt.Errorf("cannot GetUserPersonaMessages; %w", err)
	}

	msgs = append(msgs, types.Message{
		Author:  "user",
		Content: content,
	})
	reqData := chatRequest{
		Model:    "gpt-3.5-turbo",
		User:     "idk",
		Messages: msgs,
	}

	marshal, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal open ai messages; %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, openAiApiUrl, bytes.NewReader(marshal))
	if err != nil {
		return nil, fmt.Errorf("cannot create new chat request; %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send open ai request; %w", err)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot send open ai request; %w", err)
	}

	var res *openAiResponse
	err = json.Unmarshal(all, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal open ai response; %w", err)
	}

	msgs = append(msgs, res.Choices[0].Message)

	return msgs, nil
}
