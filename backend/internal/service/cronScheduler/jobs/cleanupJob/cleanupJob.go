package cleanupJob

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/errorMsg"

	"github.com/rs/zerolog/log"
)

// CleanupJob represents a job to clean up old history data from the database.
// It contains the database instance and the maximum number of days of history to keep.
type CleanupJob struct {
	db             db.Database
	MaxDaysHistory int
}

// NewCleanupJob initializes a new job for cleaning hold history data.
func NewCleanupJob(db db.Database, days int) *CleanupJob {
	return &CleanupJob{db: db, MaxDaysHistory: days}
}

// Execute is the function be used with scheduler.
func (cj *CleanupJob) Execute() {
	if err := cj.db.DeleteOldHistory(cj.MaxDaysHistory); err != nil {
		log.Error().Msgf(errorMsg.ErrorDeletingOldHistory+": %q", err)
	}
	log.Info().Msgf(errorMsg.InfoOldHistoryDataDeleted)
}
