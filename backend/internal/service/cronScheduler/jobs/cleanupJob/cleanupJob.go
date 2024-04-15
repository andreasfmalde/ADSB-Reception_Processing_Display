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

func NewCleanupJob(db db.Database, days int) *CleanupJob {
	return &CleanupJob{db: db, MaxDaysHistory: days}
}

func (cj *CleanupJob) Execute() {
	if err := cj.db.DeleteOldHistory(cj.MaxDaysHistory); err != nil {
		log.Error().Msgf(errorMsg.ErrorDeletingOldHistory+": %q", err)
	}
	log.Info().Msgf(errorMsg.InfoOldHistoryDataDeleted)
}
