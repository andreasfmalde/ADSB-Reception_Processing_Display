package service

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/models"
)

type SbsService interface {
	CreateAdsbTables() error
	InsertNewSbsData(aircraft []models.AircraftCurrentModel) error
	Cleanup() error
}

type SbsServiceImpl struct {
	DB db.Database
}

// InitSbsService initializes SbsServiceImpl struct and database connection
func InitSbsService() (*SbsServiceImpl, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	return &SbsServiceImpl{DB: dbConn}, nil
}

// CreateAdsbTables creates all tables for the database schema
func (svc *SbsServiceImpl) CreateAdsbTables() error {
	err := svc.DB.BeginTx()
	if err != nil {
		return err
	}

	defer svc.DB.Rollback()

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

	err = svc.DB.CreateAircraftHistory()
	if err != nil {
		return err
	}

	return nil
}

// InsertNewSbsData adds new SBS data to the database.
// First process the SBS stream and then add that data to the database.
func (svc *SbsServiceImpl) InsertNewSbsData(aircraft []models.AircraftCurrentModel) error {
	err := svc.DB.InsertHistoryFromCurrent()
	if err != nil {
		return err
	}

	err = svc.DB.BeginTx()
	if err != nil {
		return err
	}
	defer svc.DB.Rollback()

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
func (svc *SbsServiceImpl) Cleanup() error {
	return svc.DB.DeleteOldCurrent()
}
