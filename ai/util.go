package ai

const (
	userRole      = "user"
	assistantRole = "assistant"
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
