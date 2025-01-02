package crontab

import (
	"github.com/robfig/cron/v3"
)

func RegisterCron(f func(), spec string) error {
	_, err := c.AddFunc(spec, f)
	return err
}

var c *cron.Cron

func init() {
	c = cron.New()
}

func Start() {
	c.Start()
}
