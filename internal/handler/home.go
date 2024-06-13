package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/petr-hanzl/chatwithme/internal/views"
)

type HomeHandler struct{}

func (h HomeHandler) HomeGetHandler(c echo.Context) error {
	return views.Render(c, views.Home())
}
