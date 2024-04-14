package cleanupJob

import (
	"adsb-api/internal/db"
	"github.com/rs/zerolog/log"
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
		log.Error().Msgf("error deleting old history: %q", err.Error())
	}
	log.Info().Msgf("old SBS data deleted")
}
