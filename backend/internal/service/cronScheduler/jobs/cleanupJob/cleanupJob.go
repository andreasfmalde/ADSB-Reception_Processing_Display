package cleanupJob

import (
	"adsb-api/internal/db"
	"adsb-api/internal/utility/logger"
)

type CleanupJob struct {
	db             db.Database
	MaxDaysHistory int
}

func NewCleanupJob(db db.Database, days int) *CleanupJob {
	return &CleanupJob{db: db, MaxDaysHistory: days}
}

func (cj *CleanupJob) Execute() {
	if err := cj.db.DeleteOldHistory(cj.MaxDaysHistory); err != nil {
		logger.Error.Printf("error deleting old history: %q", err.Error())
	}
}
