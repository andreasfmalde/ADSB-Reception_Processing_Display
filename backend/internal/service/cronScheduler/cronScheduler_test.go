package cronScheduler

import (
	"adsb-api/internal/global"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	m.Run()
}

func TestNewCronScheduler(t *testing.T) {
	scheduler := NewCronScheduler()
	assert.NotNil(t, scheduler)
	assert.IsType(t, &CronScheduler{}, scheduler)
}

func TestCronScheduler_ScheduleJob(t *testing.T) {
	scheduler := NewCronScheduler()
	err := scheduler.ScheduleJob("@every 1s", func() { /* Do Nothing */ })

	assert.NoError(t, err)
	assert.Len(t, scheduler.cron.Entries(), 1)
}

func TestCronScheduler_Start(t *testing.T) {
	scheduler := NewCronScheduler()

	scheduler.Start()

	assert.NotNil(t, scheduler.cron)
}
