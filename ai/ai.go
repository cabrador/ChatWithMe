package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/petr-hanzl/chatwithme/db"
	"github.com/petr-hanzl/chatwithme/types"
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

	userMsg := types.Message{
		Author:      userRole,
		AuthorId:    db.UserAuthorId,
		Content:     content,
		UserId:      userId,
		PersonaId:   personaId,
		OrderNumber: len(msgs) + 1, // +1 because if one message is stored, this one is second
	}

	msgs = append(msgs, userMsg)
	reqData := newRequestData(msgs)

	// add rest of the data
	// this needs to be done here because OpenAi api does not allow additional data

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

	generatedMsg := res.Choices[0].Message
	generatedMsg.Author = assistantRole
	generatedMsg.AuthorId = db.AssistantAuthorId
	generatedMsg.UserId = userId
	generatedMsg.PersonaId = personaId
	generatedMsg.OrderNumber = len(msgs) + 1 // +1 because we already appended users msg

	msgsToInsert := []types.Message{userMsg, generatedMsg}
	_, err = c.db.InsertMessages(msgsToInsert)
	if err != nil {
		return nil, fmt.Errorf("cannot GetUserPersonaMessages; %w", err)
	}

	return append(msgs, generatedMsg), nil
}

func newRequestData(msgs []types.Message) chatRequest {
	var reqMsgs []openAiMessage

	for _, m := range msgs {
		reqMsgs = append(reqMsgs, openAiMessage{
			Role:    m.Author,
			Content: m.Content,
		})
	}

	return chatRequest{
		Model:    "gpt-3.5-turbo",
		User:     strconv.FormatInt(int64(msgs[0].UserId), 10),
		Messages: reqMsgs,
	}
}
