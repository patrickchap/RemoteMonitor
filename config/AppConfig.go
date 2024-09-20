package config

import (
	"github.com/robfig/cron/v3"
)

type AppConfig struct {
	Schedual    *cron.Cron
	SchedualIds map[int64]cron.EntryID
}
