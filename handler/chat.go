package handler

import (
	"chatwithme/ai"
	"chatwithme/views"
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

func NewChatHandler(generator ai.ChatGenerator) ChatHandler {
	return ChatHandler{
		generator: generator,
	}
}

type ChatHandler struct {
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

	msgs, err := h.generator.Generate(r.UserId, personaId, r.Content)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot generate chat; %w", err))
		c.Set("code", http.StatusInternalServerError)
		return echo.ErrInternalServerError
	}

	err = c.JSON(http.StatusOK, msgs)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot create json response; %w", err))
		c.Set("code", http.StatusInternalServerError)
		return echo.ErrInternalServerError
	}

	return nil
}

func (h *ChatHandler) ChatGetHandler(c echo.Context) error {
	return views.Render(c, views.Hello("world"))
}
