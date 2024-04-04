package sbsService

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/global/models"
	"adsb-api/internal/sbs"
)

type SbsService interface {
	CreateAdsbTables() error
	InsertNewSbsData(aircraft []models.AircraftCurrentModel) error
	Cleanup() error
}

type SbsImpl struct {
	DB db.Database
}

// InitSbsService initializes SbsImpl struct and database connection
func InitSbsService() (*SbsImpl, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	return &SbsImpl{DB: dbConn}, nil
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
			svc.DB.Rollback()
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

// Cleanup remove old rows to save space
func (svc *SbsImpl) Cleanup() error {
	return svc.DB.DeleteOldHistory(global.MaxDaysHistory)
}

func (svc *SbsImpl) ProcessSbsData() ([]models.AircraftCurrentModel, error) {
	return sbs.ProcessSbsStream()
}
