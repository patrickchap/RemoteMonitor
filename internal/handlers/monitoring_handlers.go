package handlers

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Job struct {
	HostServiceId int64
}

func (j *Job) Run() {
	fmt.Printf("Running job: %v", j)
	SendEvent("public", fmt.Sprintf("Running job: %v", j.HostServiceId))
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
		scheduleId, err := h.AppConfig.Schedual.AddJob(spec, &job)
		if err != nil {
			fmt.Printf("Error adding job: %v", err)
		}
		fmt.Printf("Job added: %v", scheduleId)
		h.AppConfig.SchedualIds[hostService.ID] = scheduleId

		payload := make(map[string]string)
		payload["message"] = "scheduling"
		payload["host_service_id"] = strconv.FormatInt(hostService.ID, 10)
		SendEvent("monitor", payload)

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
