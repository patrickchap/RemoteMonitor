package helpers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderTemplate(c echo.Context, component templ.Component) error {
	return echo.WrapHandler(templ.Handler(component))(c)
}
