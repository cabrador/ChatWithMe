package handler

import (
	"encoding/json"
	"fmt"
	"io"
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
		c.Set("error", fmt.Errorf("cannot convert personaId '%v' from string to int; %w", str, err))
		return echo.ErrBadRequest
	}

	msgs, err := h.ctrl.Generate(r.UserId, personaId, r.Content)
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

func (h *ChatHandler) PersonaGetHandler(c echo.Context) error {
	// todo controller

	return views.Render(c, views.Persona(c.Param("personaId")))
}

func (h *ChatHandler) ChatGetHandler(c echo.Context) error {
	return views.Render(c, views.ChatHome())
}
