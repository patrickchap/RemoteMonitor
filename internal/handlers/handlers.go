package handlers

import (
	db "RemoteMonitor/internal/database/sqlc"
	"RemoteMonitor/internal/helpers"
	"RemoteMonitor/views"
	"RemoteMonitor/views/viewModels"
	"database/sql"
	"net/http"
	"strconv"

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
	req := new(getHostParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	params := db.GetHostsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	hosts, err := h.Store.GetHosts(c.Request().Context(), params)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	var hostList []viewModels.HostList
	for _, host := range hosts {
		hostList = append(hostList, viewModels.HostList{
			HostName:      host.HostName,
			CanonicalName: nullStringToString(host.CanonicalName),
			Url:           nullStringToString(host.Url),
			Ip:            nullStringToString(host.Ip),
			Ipv6:          nullStringToString(host.Ipv6),
			Location:      nullStringToString(host.Location),
			Os:            nullStringToString(host.Os),
			Active:        nullInt64ToString(host.Active),
			LastUpdated:   nullTimetoString(host.LastUpdated),
		})
	}

	return helpers.RenderTemplate(c, views.Home(hostList))
}

func (h *Handler) Hosts(c echo.Context) error {
	req := new(getHostParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	params := db.GetHostsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	hosts, err := h.Store.GetHosts(c.Request().Context(), params)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	var hostList []viewModels.HostList
	for _, host := range hosts {
		hostList = append(hostList, viewModels.HostList{
			HostName:      host.HostName,
			CanonicalName: nullStringToString(host.CanonicalName),
			Url:           nullStringToString(host.Url),
			Ip:            nullStringToString(host.Ip),
			Ipv6:          nullStringToString(host.Ipv6),
			Location:      nullStringToString(host.Location),
			Os:            nullStringToString(host.Os),
			Active:        nullInt64ToString(host.Active),
			LastUpdated:   nullTimetoString(host.LastUpdated),
		})
	}

	return helpers.RenderTemplate(c, views.Hosts(hostList))
}
func (h *Handler) Login(c echo.Context) error {
	var isLoggedIn bool
	isLoggedIn = false
	if isLoggedIn {
		return c.Redirect(302, "/dashboard")
	}

	return helpers.RenderTemplate(c, views.Login())
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
