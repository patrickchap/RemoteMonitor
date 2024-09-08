package handlers

import (
	db "RemoteMonitor/internal/database/sqlc"
	"RemoteMonitor/internal/helpers"
	"RemoteMonitor/views"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Store   db.Store
	Manager *manager
}

type manager struct {
	session    *sessions.Session
	cookie     CookieOpts
	authFailed echo.HandlerFunc
}

type CookieOpts struct {
	Name   string
	Secret string
	MaxAge int
}

type getHostParams struct {
	Offset int64 `query:"page"`
	Limit  int64 `query:"limit"`
}

func (h *Handler) Dashboard(c echo.Context) error {

	u, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	fmt.Printf(">>>>>>> user: %v", u)

	req := new(getHostParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	params := db.GetHostsWithServicesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	hosts, err := h.Store.GetHostsWithServices(c.Request().Context(), params)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return helpers.RenderTemplate(c, views.Home(hosts))
}

func (h *Handler) Hosts(c echo.Context) error {
	req := new(getHostParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	params := db.GetHostsWithServicesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	hosts, err := h.Store.GetHostsWithServices(c.Request().Context(), params)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return helpers.RenderTemplate(c, views.Hosts(hosts))
}

type HostEditParams struct {
	Id int64 `param:"id"`
}

func (h *Handler) HostEdit(c echo.Context) error {
	req := new(HostEditParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	fmt.Printf("request: %d", req)
	return helpers.RenderTemplate(c, views.HostEdit(req.Id))
}

func (h *Handler) WsTest(c echo.Context) error {
	return helpers.RenderTemplate(c, views.WebsocketClient())
}

func nullInt64ToString(i sql.NullInt64) string {
	if i.Valid {
		return strconv.FormatInt(i.Int64, 10)
	}
	return ""
}

func nullTimetoString(i sql.NullTime) string {
	if i.Valid {
		return i.Time.String()
	}
	return ""
}

func nullStringToString(i sql.NullString) string {
	if i.Valid {
		return i.String
	}
	return ""
}
