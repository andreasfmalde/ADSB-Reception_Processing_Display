package service

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/models"
	"adsb-api/internal/sbs"
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
	err := svc.DB.CreateAircraftCurrentTable(nil)
	if err != nil {
		return err
	}

	err = svc.DB.CreateAircraftHistory(nil)
	if err != nil {
		return err
	}

	return nil
}

// InsertNewSbsData adds new SBS data to the database.
// First process the SBS stream and then add that data to the database.
func (svc *SbsServiceImpl) InsertNewSbsData() error {
	err := svc.DB.InsertHistoryFromCurrent(nil)
	if err != nil {
		return err
	}

	aircraft, err := sbs.ProcessSBSstream()
	if err != nil {
		return err
	}

	tx, err := svc.DB.Begin()

	err = svc.DB.DropAircraftCurrentTable(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = svc.DB.CreateAircraftCurrentTable(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = svc.DB.BulkInsertAircraftCurrent(aircraft, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = svc.DB.Commit(tx)
	if err != nil {
		return err
	}

	return nil
}

// Cleanup remove old rows to save space
func (svc *SbsServiceImpl) Cleanup() error {
	return svc.DB.DeleteOldCurrent(nil)
}
