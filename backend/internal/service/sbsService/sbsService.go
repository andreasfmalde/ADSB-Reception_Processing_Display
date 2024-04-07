package sbsService

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/models"
	"adsb-api/internal/sbs"
	"adsb-api/internal/service/cronScheduler"
	"adsb-api/internal/utility/logger"
	"errors"
)

type SbsService interface {
	CreateAdsbTables() error
	InsertNewSbsData(aircraft []models.AircraftCurrentModel) error
	InitAndStartCleanUpJob() error
	StartScheduler() error
}

type SbsImpl struct {
	DB            db.Database
	CronScheduler cronScheduler.Scheduler
}

// InitSbsService initializes SbsImpl struct and database connection
func InitSbsService() (*SbsImpl, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	scheduler := cronScheduler.NewCronScheduler()

	return &SbsImpl{DB: dbConn, CronScheduler: scheduler}, nil
}

// CreateAdsbTables creates all tables for the database schema
func (svc *SbsImpl) CreateAdsbTables() error {
	err := svc.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			svc.DB.Rollback()
		}
	}()

	err = svc.DB.CreateAircraftCurrentTable()
	if err != nil {
		return err
	}

	err = svc.DB.CreateAircraftCurrentTimestampIndex()
	if err != nil {
		return err
	}

	err = svc.DB.Commit()
	if err != nil {
		return err
	}

	err = svc.DB.CreateAircraftHistoryTable()
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
			svc.DB.Rollback()
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

func (svc *SbsImpl) ProcessSbsData() ([]models.AircraftCurrentModel, error) {
	return sbs.ProcessSbsStream()
}

func (svc *SbsImpl) StartScheduler() error {
	if svc.CronScheduler == nil {
		return errors.New(errorMsg.CronSchedulerIsNotInitialized)
	}
	svc.CronScheduler.Start()
	return nil
}

func (svc *SbsImpl) CleanupJob() {
	err := svc.DB.DeleteOldHistory(global.MaxDaysHistory)
	if err != nil {
		logger.Error.Printf("error deleting old history: %q", err.Error())
	}
}

// ScheduleCleanUpJob initializes a starts a cleanup job, that remove old rows from database to save space.
// When the job is scheduled to be executed is decided by schedule parameter.
func (svc *SbsImpl) ScheduleCleanUpJob(schedule string) error {
	if svc.CronScheduler == nil {
		return errors.New(errorMsg.CronSchedulerIsNotInitialized)
	}
	return svc.CronScheduler.ScheduleJob(schedule, svc.CleanupJob)
}
