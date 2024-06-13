package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/petr-hanzl/chatwithme/internal/types"
)

const (
	openAiApiUrl      = "https://api.openai.com/v1/chat/completions"
	UserRole          = "user"
	AssistantRole     = "assistant"
	UserAuthorId      = 1
	AssistantAuthorId = 2
)

type chatRequest struct {
	Model    string          `json:"model"`
	User     string          `json:"user"`
	Messages []openAiMessage `json:"messages"`
}

type openAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

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

func MakeGenerator() *Generator {
	return &Generator{
		client: http.DefaultClient,
		apiKey: os.Getenv("OPENAI_API_TOKEN"),
	}
}

type Generator struct {
	client *http.Client
	apiKey string
}

func (g *Generator) Generate(msgs []types.Message, userId, personaId int) (types.Message, error) {
	reqData := newRequestData(msgs)

	// add rest of the data
	// this needs to be done here because OpenAi api does not allow additional data

	marshal, err := json.Marshal(reqData)
	if err != nil {
		return types.Message{}, fmt.Errorf("cannot marshal open ai messages; %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, openAiApiUrl, bytes.NewReader(marshal))
	if err != nil {
		return types.Message{}, fmt.Errorf("cannot create new chat request; %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", g.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return types.Message{}, fmt.Errorf("cannot send open ai request; %w", err)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Message{}, fmt.Errorf("cannot send open ai request; %w", err)
	}

	var res *openAiResponse
	err = json.Unmarshal(all, &res)
	if err != nil {
		return types.Message{}, fmt.Errorf("cannot unmarshal open ai response; %w", err)
	}

	generatedMsg := res.Choices[0].Message
	generatedMsg.Author = AssistantRole
	generatedMsg.AuthorId = AssistantAuthorId
	generatedMsg.UserId = userId
	generatedMsg.PersonaId = personaId
	generatedMsg.OrderNumber = len(msgs) + 1 // +1 because we already appended users msg

	return generatedMsg, nil
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
