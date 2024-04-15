package cronScheduler

import (
	"github.com/robfig/cron"
)

// Scheduler is an interface for scheduling jobs based on a given schedule.
type Scheduler interface {
	ScheduleJob(schedule string, job func()) error
	Start()
}

// CronScheduler is a type that represents a scheduler based on the cron library.
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
