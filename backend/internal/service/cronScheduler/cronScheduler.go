package cronScheduler

import (
	"github.com/robfig/cron"
)

type Scheduler interface {
	ScheduleJob(schedule string, job func()) error
	Start()
}

type CronScheduler struct {
	cron *cron.Cron
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron: cron.New(),
	}
}

func (cs *CronScheduler) ScheduleJob(schedule string, job func()) error {
	return cs.cron.AddFunc(schedule, job)
}

func (cs *CronScheduler) Start() {
	cs.cron.Start()
}
