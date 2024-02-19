package ai

type openAiRole string

const (
	userRole      openAiRole = "user"
	assistantRole openAiRole = "assistant"
)

type chatMessage struct {
	Role    openAiRole
	Content string
}

type chatRequest struct {
	Model    string
	User     string
	Messages []chatMessage
}
