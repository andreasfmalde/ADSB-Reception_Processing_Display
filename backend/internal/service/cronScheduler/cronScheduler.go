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

// NewCronScheduler creates a new instance of CronScheduler,
// initializing it with a new cron object.
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron: cron.New(),
	}
}

// ScheduleJob schedules a job to CronScheduler cron object.
func (cs *CronScheduler) ScheduleJob(schedule string, job func()) error {
	return cs.cron.AddFunc(schedule, job)
}

// Start starts cron scheduler.
// Every job schedules before and after the call of this method will be executed.
func (cs *CronScheduler) Start() {
	cs.cron.Start()
}
