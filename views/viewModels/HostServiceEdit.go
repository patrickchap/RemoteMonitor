package viewmodels

type HostServiceEdit struct {
	Id             int64  `json:"id"`
	HostId         int64  `json:"host_id"`
	HostName       string `json:"host_name"`
	ServiceId      int64  `json:"service_id"`
	ServiceName    string `json:"service_name"`
	Active         int64  `json:"active"`
	ScheduleNumber int64  `json:"schedule_number"`
	ScheduleUnit   string `json:"schedule_unit"`
}

type Service struct {
	ServiceId   int64  `json:"service_id"`
	ServiceName string `json:"service_name"`
}
