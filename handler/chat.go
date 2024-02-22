package handler

import (
	"chatwithme/ai"
	"chatwithme/db"
	"chatwithme/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type messageReq struct {
	UserId      int
	Content     string
	OrderNumber int
}

func NewChatHandler(db *db.Database, generator ai.ChatGenerator) ChatHandler {
	return ChatHandler{
		db:        db,
		generator: generator,
	}
}

type ChatHandler struct {
	db        *db.Database
	generator ai.ChatGenerator
}

func (h *ChatHandler) ChatPostHandler(c echo.Context) error {
	req := c.Request()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot read body; %w", err))
		return echo.ErrBadRequest
	}

	var r messageReq
	err = json.Unmarshal(data, &r)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot unmarshal message req; %w", err))
		return echo.ErrBadRequest
	}
	str := c.Param("personaId")
	personaId, err := strconv.Atoi(str)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot convert personaId '%v from string to int; %w", str, err))
		return echo.ErrBadRequest
	}

	msgs, err := h.db.GetUserPersonaMessages(r.UserId, personaId)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot GetUserPersonaMessages; %w", err))
		return echo.ErrInternalServerError
	}

	msgs = append(msgs, types.Message{
		Author:  "user",
		Content: r.Content,
	})

	err = h.generator.Generate(msgs)
	if err != nil {
		return err
	}

	//_, err = h.db.InsertMessage(r.UserId, personaId, r.Content, r.OrderNumber)
	//if err != nil {
	//	c.Set("error", fmt.Errorf("cannot GetUserPersonaMessages; %w", err))
	//	return echo.ErrInternalServerError
	//}

	return nil
}

func (h *ChatHandler) ChatGetHandler(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "not yet implemented", http.StatusServiceUnavailable)
}
