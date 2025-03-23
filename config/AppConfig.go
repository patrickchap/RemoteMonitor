package config

import (
	"github.com/robfig/cron/v3"
)

type AppConfig struct {
	Schedual      *cron.Cron
	SchedualIds   map[int64]cron.EntryID
	ShouldMonitor bool
}

func (c *AppConfig) SetShouldMonitor(shouldMonitor bool) {
	c.ShouldMonitor = shouldMonitor
}

func (c *AppConfig) GetShouldMonitor() bool {
	return c.ShouldMonitor
}
