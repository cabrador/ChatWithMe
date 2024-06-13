package db

import (
	"fmt"
	"net/http"

	"github.com/petr-hanzl/chatwithme/internal/ai"
	"github.com/petr-hanzl/chatwithme/internal/types"
)

func MakeChatController(db *Database) ChatController {
	return &chatController{
		db:        db,
		generator: ai.MakeGenerator(),
	}
}

type ChatController interface {
	// Generate sends the given message to OpenAI chat.
	// If there were any messages sent before, they
	// are retrieved from the database.
	Generate(userId, personaId int, content string) (types.Messages, error)

	GetUserPersonaMessages(userId, personaId int) ([]types.Message, error)
}

type chatController struct {
	client    *http.Client
	apiKey    string
	db        *Database
	generator *ai.Generator
}

func (c *chatController) GetUserPersonaMessages(userId, personaId int) ([]types.Message, error) {
	return c.db.GetUserPersonaMessages(userId, personaId)
}

func (c *chatController) Generate(userId, personaId int, content string) (types.Messages, error) {
	msgs, err := c.db.GetUserPersonaMessages(userId, personaId)
	if err != nil {
		return nil, fmt.Errorf("cannot GetUserPersonaMessages; %w", err)
	}

	userMsg := types.Message{
		Author:      ai.UserRole,
		AuthorId:    ai.UserAuthorId,
		Content:     content,
		UserId:      userId,
		PersonaId:   personaId,
		OrderNumber: len(msgs) + 1, // +1 because if one message is stored, this one is second
	}
	aiAnswer, err := c.generator.Generate(msgs, userId, personaId)
	if err != nil {
		return nil, err
	}

	msgsToInsert := []types.Message{userMsg, aiAnswer}
	_, err = c.db.InsertMessages(msgsToInsert)
	if err != nil {
		return nil, fmt.Errorf("cannot GetUserPersonaMessages; %w", err)
	}

	return append(msgs, msgsToInsert...), nil
}
