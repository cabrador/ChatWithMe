package ai

import (
	"bytes"
	"chatwithme/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	openAiApiUrl = "https://api.openai.com/v1/chat/completions"
)

func MakeChatGenerator() ChatGenerator {
	return &chatGenerator{
		client: http.DefaultClient,
		apiKey: os.Getenv("OPENAI_API_TOKEN"),
	}
}

type ChatGenerator interface {
	Generate([]db.Message) error
}

type chatGenerator struct {
	client *http.Client
	apiKey string
}

func (c *chatGenerator) Generate(msgs []db.Message) error {
	reqData := chatRequest{
		Model:    "gpt-3.5-turbo",
		User:     "idk",
		Messages: msgs,
	}

	marshal, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("cannot marshal open ai messages; %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, openAiApiUrl, bytes.NewReader(marshal))
	if err != nil {
		return fmt.Errorf("cannot create new chat request; %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send open ai request; %w", err)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(all)

	return nil
}
