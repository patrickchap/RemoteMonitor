package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Job struct {
	HostServiceId int64
	Handler       *Handler
}

func (j *Job) Run() {
	fmt.Printf("Running job: %v\n", j)
	j.Handler.CheckHostService(j.HostServiceId)
}

func (h *Handler) CheckHostService(hostServiceId int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	hostService, err := h.Store.GetHostService(ctx, hostServiceId)
	if err != nil {
		fmt.Printf("Error getting host service: %v", err)
	}

	var message, status string
	switch hostService.ServiceName.String {
	case "http":
		message, status = CheckHttpService(hostService.Url.String)
	}

	fmt.Printf("Message: %v\n", message)
	fmt.Printf("Status: %v\n", status)
	fmt.Printf("Checking host service: %v\n", hostService)
	//TODO: Send  message
	//SendEvent("host-channel", message)
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

//TODO: Stop monitorying

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
