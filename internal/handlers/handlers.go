package handlers

import (
	db "RemoteMonitor/internal/database/sqlc"
	"RemoteMonitor/internal/helpers"
	"RemoteMonitor/views"
	viewmodels "RemoteMonitor/views/viewModels"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	return helpers.RenderTemplate(c, views.HostEdit(req.Id))
}

type HostEditFormParams struct {
	Id int64 `param:"id"`
}

func (h *Handler) GetEditHostDetails(c echo.Context) error {
	req := new(HostEditFormParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	host, err := h.Store.GetHost(c.Request().Context(), req.Id)
	if err != nil {

		return c.String(http.StatusBadRequest, "Bad Request")
	}

	viewHost := viewmodels.Host{
		ID:            host.ID,
		HostName:      host.HostName,
		CanonicalName: host.CanonicalName,
		Url:           host.Url,
		Ip:            host.Ip,
		Ipv6:          host.Ipv6,
	}
	return helpers.RenderTemplate(c, views.EditHostForm(viewHost))

}

type PutEditHostParams struct {
	ID            int64  `form:"id"`
	HostName      string `form:"host_name"`
	CanonicalName string `form:"canonical_name"`
	Url           string `form:"url"`
	Ip            string `form:"ip"`
	Ipv6          string `form:"ipv6"`
}

func (h *Handler) PutEditHostDetails(c echo.Context) error {
	req := new(PutEditHostParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	updateHost := db.UpdateHostParams{
		ID:            req.ID,
		HostName:      req.HostName,
		CanonicalName: sql.NullString{String: req.CanonicalName, Valid: true},
		Url:           sql.NullString{String: req.Url, Valid: true},
		Ip:            sql.NullString{String: req.Ip, Valid: true},
		Ipv6:          sql.NullString{String: req.Ipv6, Valid: true},
		LastUpdated:   sql.NullTime{Time: time.Now(), Valid: true},
	}

	host, err := h.Store.UpdateHost(c.Request().Context(), updateHost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	viewHost := viewmodels.Host{
		HostName:      host.HostName,
		CanonicalName: host.CanonicalName,
		Url:           host.Url,
		Ip:            host.Ip,
		Ipv6:          host.Ipv6,
	}

	return helpers.RenderTemplate(c, views.EditHostForm(viewHost))
}

func (h *Handler) HostCreateForm(c echo.Context) error {
	return helpers.RenderTemplate(c, views.CreateHostForm())
}

type PostCreateHostParams struct {
	HostName      string `form:"host_name"`
	CanonicalName string `form:"canonical_name"`
	Url           string `form:"url"`
	Ip            string `form:"ip"`
	Ipv6          string `form:"ipv6"`
}

func (h *Handler) HostCreate(c echo.Context) error {
	req := new(PostCreateHostParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	params := db.CreateHostParams{
		HostName:      req.HostName,
		CanonicalName: sql.NullString{String: req.CanonicalName, Valid: true},
		Url:           sql.NullString{String: req.Url, Valid: true},
		Ip:            sql.NullString{String: req.Ip, Valid: true},
		Ipv6:          sql.NullString{String: req.Ipv6, Valid: true},
	}

	host, err := h.Store.CreateHost(c.Request().Context(), params)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/admin/host/edit/%d", host.ID))
	return nil
}

type GetHostServicesParams struct {
	HostId int64 `param:"host_id"`
}

func (h *Handler) GetHostServices(c echo.Context) error {
	req := new(GetHostServicesParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	hostServices, err := h.Store.GetHostServices(c.Request().Context(), sql.NullInt64{Int64: req.HostId, Valid: true})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	acitveService, err := h.Store.GetServices(c.Request().Context())

	availableServices := []viewmodels.Service{}
	for _, service := range acitveService {
		found := false
		for _, hostService := range hostServices {
			if service.ID == hostService.ServiceID.Int64 {
				found = true
				break
			}
		}
		if !found {
			availableServices = append(availableServices, viewmodels.Service{
				ServiceId:   service.ID,
				ServiceName: service.ServiceName.String,
			})
		}
	}

	fmt.Printf(">>>>> availableService: %v", availableServices)
	hostServiceModel := []viewmodels.HostServiceEdit{}

	for _, hostService := range hostServices {
		hostServiceModel = append(hostServiceModel, viewmodels.HostServiceEdit{
			HostId:         hostService.HostID.Int64,
			HostName:       hostService.HostName,
			ServiceId:      hostService.ServiceID.Int64,
			ServiceName:    hostService.ServiceName.String,
			Active:         hostService.Active.Int64,
			ScheduleNumber: hostService.ScheduleNumber.Int64,
			ScheduleUnit:   hostService.ScheduleUnit.String,
		})
	}

	return helpers.RenderTemplate(c, views.EditServicesForm(hostServiceModel, availableServices))
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
