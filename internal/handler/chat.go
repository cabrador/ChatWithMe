package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/petr-hanzl/chatwithme/internal/db"
)

type MessageReq struct {
	UserId  int
	Content string
}

func NewChatHandler(controller db.ChatController) ChatHandler {
	return ChatHandler{
		ctrl: controller,
	}
}

type ChatHandler struct {
	ctrl db.ChatController
}

func (h *ChatHandler) PersonaPostHandler(c echo.Context) error {
	// todo get userid from context
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot read request body; %w", err))
		c.Set("code", http.StatusInternalServerError)
		return echo.ErrInternalServerError
	}

	var req MessageReq
	err = json.Unmarshal(body, &req)
	if err != nil || req.UserId == 0 || req.Content == "" {
		c.Set("error", fmt.Errorf("incorrect request body; %w", err))
		return echo.ErrBadRequest
	}

	strPID := c.Param("personaId")
	personaId, err := strconv.Atoi(strPID)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot convert personaId '%v' from string to int; %w", strPID, err))
		return echo.ErrBadRequest
	}

	msgs, err := h.ctrl.Generate(req.UserId, personaId, req.Content)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot generate chat; %w", err))
		c.Set("code", http.StatusInternalServerError)
		return echo.ErrInternalServerError
	}

	err = c.String(http.StatusOK, msgs.String())
	if err != nil {
		c.Set("error", fmt.Errorf("cannot create json response; %w", err))
		c.Set("code", http.StatusInternalServerError)
		return echo.ErrInternalServerError
	}

	return nil
}
