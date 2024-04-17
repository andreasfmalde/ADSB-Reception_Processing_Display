package sbsService

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/models"
	"adsb-api/internal/service/cronScheduler"
	"adsb-api/internal/service/cronScheduler/jobs/cleanupJob"
)

// SbsService represents a service with an interface for retrieving database data through the repository in
// internal/db/database.go.
type SbsService interface {
	CreateAdsbTables() error
	InsertNewSbsData(aircraft []models.AircraftCurrentModel) error
	ScheduleCleanUpJob(schedule string, days int) error
}

type SbsImpl struct {
	DB            db.Database
	CronScheduler cronScheduler.Scheduler
}

// InitSbsService initializes SbsImpl struct and database connection
func InitSbsService(db db.Database, scheduler cronScheduler.Scheduler) *SbsImpl {
	return &SbsImpl{DB: db, CronScheduler: scheduler}
}

// CreateAdsbTables creates all tables for the database schema
func (svc *SbsImpl) CreateAdsbTables() error {
	err := svc.DB.CreateAircraftCurrentTable()
	if err != nil {
		return err
	}

	err = svc.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = svc.DB.Rollback()
		}
	}()

	err = svc.DB.CreateAircraftHistoryTable()
	if err != nil {
		return err
	}

	err = svc.DB.CreateAircraftHistoryTimestampIndex()
	if err != nil {
		return err
	}

	err = svc.DB.Commit()
	if err != nil {
		return err
	}

	return nil
}

// InsertNewSbsData adds new SBS data to the database.
// First process the SBS stream and then add that data to the database.
func (svc *SbsImpl) InsertNewSbsData(aircraft []models.AircraftCurrentModel) error {
	err := svc.DB.InsertHistoryFromCurrent()
	if err != nil {
		return err
	}

	err = svc.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = svc.DB.Rollback()
		}
	}()

	err = svc.DB.DropAircraftCurrentTable()
	if err != nil {
		return err
	}

	err = svc.DB.CreateAircraftCurrentTable()
	if err != nil {
		return err
	}

	err = svc.DB.BulkInsertAircraftCurrent(aircraft)
	if err != nil {
		return err
	}

	err = svc.DB.Commit()
	if err != nil {
		return err
	}

	return nil
}

// ScheduleCleanUpJob initializes a starts a cleanupJob job, that remove old rows from database to save space.
// When the job is scheduled to be executed is decided by schedule parameter.
func (svc *SbsImpl) ScheduleCleanUpJob(schedule string, days int) error {
	job := cleanupJob.NewCleanupJob(svc.DB, days)
	return svc.CronScheduler.ScheduleJob(schedule, job.Execute)
}

// StartScheduler starts the cron scheduler.
// Every job scheduled before this method is called will begin.
// Jobs scheduled after this method will still be executed.
func (svc *SbsImpl) StartScheduler() {
	svc.CronScheduler.Start()
}
