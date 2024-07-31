package handlers

import (
	"RemoteMonitor/internal/helpers"
	"RemoteMonitor/views"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func (h *Handler) Dashboard(c echo.Context) error {
	return helpers.RenderTemplate(c, views.Home())
}
