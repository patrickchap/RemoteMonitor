package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	database "RemoteMonitor/internal/database/sqlc"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Job struct {
	HostServiceId int64
	Handler       *Handler
}

func (j *Job) Run() {
	fmt.Printf("Running job: %v\n", j)
	j.Handler.CheckHostService(j.HostServiceId)
}

type HostServiceStatusPayload struct {
	HostServiceId int64  `json:"host_service_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	LastCheck     string `json:"last_check"`
}

func (h *Handler) CheckHostService(hostServiceId int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	hostService, err := h.Store.GetHostService(ctx, hostServiceId)
	if err != nil {
		fmt.Printf("Error getting host service: %v", err)
		return
	}

	var message, status string
	switch hostService.ServiceName.String {
	case "http":
		message, status = CheckHttpService(hostService.Url.String)
	}

	fmt.Printf("Message: %v\n", message)
	fmt.Printf("Status: %v\n", status)
	fmt.Printf("Checking host service: %v\n", hostService)

	// Persist the status to the database
	updateCtx, updateCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer updateCancel()
	_, err = h.Store.UpdateHostServiceStatus(updateCtx, database.UpdateHostServiceStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     hostServiceId,
	})
	if err != nil {
		fmt.Printf("Error updating host service status: %v\n", err)
	}

	// Send status update to all connected clients
	payload := HostServiceStatusPayload{
		HostServiceId: hostServiceId,
		Status:        status,
		Message:       message,
		LastCheck:     time.Now().Format(time.RFC3339),
	}
	SendEvent("host-service-status", payload)
}

// TODO: update return type in order to send event types to different channels
func CheckHttpService(url string) (string, string) {
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	url = strings.Replace(url, "https://", "http://", -1)

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	fmt.Printf("Checking url: %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}

// TODO: Stop monitorying
func (h *Handler) StopMonitor() {
	fmt.Println("Stopping monitor")

}

func (h *Handler) Monitor() {

	fmt.Println("Monitoring started")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	hostServices, err := h.Store.GetHostsServicesToMonitor(ctx)
	if err != nil {
		fmt.Printf("Error getting host services: %v", err)
	}

	for _, hostService := range hostServices {

		spec := fmt.Sprintf("@every %d%s", hostService.ScheduleNumber.Int64, hostService.ScheduleUnit.String)

		var job Job
		job.HostServiceId = hostService.ID
		job.Handler = h
		scheduleId, err := h.AppConfig.Schedual.AddJob(spec, &job)
		if err != nil {
			fmt.Printf("Error adding job: %v", err)
		}
		fmt.Printf("Job added: %v", scheduleId)
		h.AppConfig.SchedualIds[hostService.ID] = scheduleId

		payload := make(map[string]string)
		payload["message"] = "scheduling"
		payload["host_service_id"] = strconv.FormatInt(hostService.ID, 10)

	}

}

type Event struct {
	Name string
	Data interface{}
}

var EventChannel = make(chan Event)

var (
	WsClients = make(map[*websocket.Conn]bool)
	WsMutex   sync.Mutex
)

func SendEvent(name string, data interface{}) {
	EventChannel <- Event{Name: name, Data: data}
}

func (h *Handler) ToggleMonitor(c echo.Context) error {
	currentState := !h.AppConfig.GetShouldMonitor()
	h.AppConfig.SetShouldMonitor(currentState)
	fmt.Printf("Current state: %v\n", currentState)
	if currentState {
		// Starting monitoring
		fmt.Println("Starting monitoring toogle")
		h.AppConfig.Schedual.Start()
	} else {
		// Stopping monitoring
		fmt.Println("Stopping monitoring toogle")
		h.AppConfig.Schedual.Stop()
	}

	return c.JSON(http.StatusOK, map[string]bool{"monitoring": currentState})
}

func (h *Handler) GetMonitorState(c echo.Context) error {
	currentState := h.AppConfig.GetShouldMonitor()
	fmt.Printf("Current state: %v\n", currentState)
	return c.JSON(http.StatusOK, map[string]bool{"monitoring": currentState})
}
