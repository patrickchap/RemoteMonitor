package handlers

import (
	db "RemoteMonitor/internal/database/sqlc"
	"RemoteMonitor/internal/helpers"
	"RemoteMonitor/views"
	component "RemoteMonitor/views/components"
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
	HostId int64 `param:"id"`
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

	hostServiceModel := []viewmodels.HostServiceEdit{}

	for _, hostService := range hostServices {
		hostServiceModel = append(hostServiceModel, viewmodels.HostServiceEdit{
			Id:             hostService.ID,
			HostId:         hostService.HostID.Int64,
			HostName:       hostService.HostName,
			ServiceId:      hostService.ServiceID.Int64,
			ServiceName:    hostService.ServiceName.String,
			Active:         hostService.Active.Int64,
			ScheduleNumber: hostService.ScheduleNumber.Int64,
			ScheduleUnit:   hostService.ScheduleUnit.String,
		})
	}
	fmt.Println("EditServicesForm")
	fmt.Println(req.HostId)
	scripts := views.EmptyScripts()
	return helpers.RenderTemplate(c, views.EditServicesForm(hostServiceModel, availableServices, req.HostId, scripts))
}

type PostHostServiceParams struct {
	HostId    int64 `form:"host_id"`
	ServiceId int64 `form:"service_id"`
}

func (h *Handler) PostHostService(c echo.Context) error {
	req := new(PostHostServiceParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	params := db.CreateHostServiceParams{
		HostID:    sql.NullInt64{Int64: req.HostId, Valid: true},
		ServiceID: sql.NullInt64{Int64: req.ServiceId, Valid: true},
	}

	innactiveParams := db.GetInnactiveHostServiceByHostAndServiceParams{
		HostID:    sql.NullInt64{Int64: req.HostId, Valid: true},
		ServiceID: sql.NullInt64{Int64: req.ServiceId, Valid: true},
	}

	innactiveHostService, err := h.Store.GetInnactiveHostServiceByHostAndService(c.Request().Context(), innactiveParams)
	if err != nil {
		_, err = h.Store.CreateHostService(c.Request().Context(), params)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	} else {
		fmt.Printf("innactive host service: %v", innactiveHostService)
		_, err = h.Store.ReactivateHostService(c.Request().Context(), innactiveHostService.ID)
	}

	fmt.Printf("innactive host service: %v", innactiveHostService)

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

	hostServiceModel := []viewmodels.HostServiceEdit{}

	for _, hostService := range hostServices {
		hostServiceModel = append(hostServiceModel, viewmodels.HostServiceEdit{
			Id:             hostService.ID,
			HostId:         hostService.HostID.Int64,
			HostName:       hostService.HostName,
			ServiceId:      hostService.ServiceID.Int64,
			ServiceName:    hostService.ServiceName.String,
			Active:         hostService.Active.Int64,
			ScheduleNumber: hostService.ScheduleNumber.Int64,
			ScheduleUnit:   hostService.ScheduleUnit.String,
		})
	}
	fmt.Println("EditServicesForm")
	fmt.Println(req.HostId)

	scripts := views.EmptyScripts()
	return helpers.RenderTemplate(c, views.EditServicesForm(hostServiceModel, availableServices, req.HostId, scripts))
}

type EditServiceRowParams struct {
	ServiceId int64 `param:"id"`
}

func (h *Handler) EditServiceRow(c echo.Context) error {
	req := new(EditServiceRowParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	hostService, err := h.Store.GetHostService(c.Request().Context(), req.ServiceId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	hostServiceModel := viewmodels.HostServiceEdit{
		Id:             hostService.ID,
		HostId:         hostService.HostID.Int64,
		HostName:       hostService.HostName,
		ServiceId:      hostService.ServiceID.Int64,
		ServiceName:    hostService.ServiceName.String,
		Active:         hostService.Active.Int64,
		ScheduleNumber: hostService.ScheduleNumber.Int64,
		ScheduleUnit:   hostService.ScheduleUnit.String,
	}

	return helpers.RenderTemplate(c, component.EditServiceRow(hostServiceModel))
}

type GetServiceRowParams struct {
	ServiceId int64 `param:"id"`
}

func (h *Handler) GetServiceRow(c echo.Context) error {
	req := new(GetServiceRowParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	hostService, err := h.Store.GetHostService(c.Request().Context(), req.ServiceId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	hostServiceModel := viewmodels.HostServiceEdit{
		Id:             hostService.ID,
		HostId:         hostService.HostID.Int64,
		HostName:       hostService.HostName,
		ServiceId:      hostService.ServiceID.Int64,
		ServiceName:    hostService.ServiceName.String,
		Active:         hostService.Active.Int64,
		ScheduleNumber: hostService.ScheduleNumber.Int64,
		ScheduleUnit:   hostService.ScheduleUnit.String,
	}

	return helpers.RenderTemplate(c, component.ServiceRow(hostServiceModel))
}

type GetDeleteServiceRowParams struct {
	ServiceId int64 `param:"id"`
}

func (h *Handler) DeleteServiceRow(c echo.Context) error {
	req := new(GetDeleteServiceRowParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	deleted, _ := h.Store.DeleteHostService(c.Request().Context(), req.ServiceId)

	hostServices, err := h.Store.GetHostServices(c.Request().Context(), deleted.HostID)
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

	hostServiceModel := []viewmodels.HostServiceEdit{}

	for _, hostService := range hostServices {
		hostServiceModel = append(hostServiceModel, viewmodels.HostServiceEdit{
			Id:             hostService.ID,
			HostId:         hostService.HostID.Int64,
			HostName:       hostService.HostName,
			ServiceId:      hostService.ServiceID.Int64,
			ServiceName:    hostService.ServiceName.String,
			Active:         hostService.Active.Int64,
			ScheduleNumber: hostService.ScheduleNumber.Int64,
			ScheduleUnit:   hostService.ScheduleUnit.String,
		})
	}

	scripts := views.DeleteSuccessfullScirpt()
	return helpers.RenderTemplate(c, views.EditServicesForm(hostServiceModel, availableServices, deleted.HostID.Int64, scripts))
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
