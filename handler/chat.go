package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/petr-hanzl/chatwithme/db"
	"github.com/petr-hanzl/chatwithme/views"
)

type messageReq struct {
	UserId      int
	Content     string
	OrderNumber int
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
	userId := 1
	content := c.FormValue("content")
	strPID := c.Param("personaId")
	personaId, err := strconv.Atoi(strPID)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot convert personaId '%v' from string to int; %w", strPID, err))
		return echo.ErrBadRequest
	}

	msgs, err := h.ctrl.Generate(userId, personaId, content)
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

func (h *ChatHandler) PersonaGetHandler(c echo.Context) error {
	// todo get userid from context
	str := c.Param("personaId")
	personaId, err := strconv.Atoi(str)
	if err != nil {
		c.Set("error", fmt.Errorf("cannot convert personaId '%v' from string to int; %w", str, err))
		return echo.ErrBadRequest
	}

	msgs, err := h.ctrl.GetUserPersonaMessages(1, personaId)
	if err != nil {
		c.Set("error", err)
		c.Set("code", http.StatusInternalServerError)
		return echo.ErrInternalServerError
	}

	return views.Render(c, views.Persona(c.Param("personaId"), msgs))
}

func (h *ChatHandler) ChatGetHandler(c echo.Context) error {
	return views.Render(c, views.ChatHome())
}
