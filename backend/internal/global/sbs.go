package global

const (
	WaitingTime    = 4
	UpdatingPeriod = 10
	MaxDaysHistory = 1
)

var (
	SbsSource       string
	CleanupSchedule = "0 0 * * *" // once a day
)
