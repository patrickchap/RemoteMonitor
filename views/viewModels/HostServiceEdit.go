package viewmodels

type HostServiceEdit struct {
	HostId         int64  `json:"host_id"`
	HostName       string `json:"host_name"`
	ServiceId      int64  `json:"service_id"`
	ServiceName    string `json:"service_name"`
	Active         bool   `json:"active"`
	ScheduleNumber int    `json:"schedule_number"`
	ScheduleUnit   string `json:"schedule_unit"`
}
