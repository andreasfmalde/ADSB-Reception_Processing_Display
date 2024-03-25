package service

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/models"
)

type SbsService interface {
	CreateAdsbTables() error
	UpdateHistory() error
	InsertNewAircraft(aircraft []models.AircraftCurrentModel) error
	Cleanup() error
}

type SbsServiceImpl struct {
	DB db.Database
}

// InitSbsService initializes SbsServiceImpl struct and database connection
func InitSbsService() (*RestServiceImpl, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	return &RestServiceImpl{DB: dbConn}, nil
}

// CreateAdsbTables creates all tables for the database schema
func (svc *SbsServiceImpl) CreateAdsbTables() error {
	return svc.DB.CreateAdsbTables()
}

// InsertNewAircraft adds new aircraft to the database
func (svc *SbsServiceImpl) InsertNewAircraft(aircraft []models.AircraftCurrentModel) error {
	return svc.DB.BulkInsertAircraftCurrent(aircraft)
}

// UpdateHistory updates history table
func (svc *SbsServiceImpl) UpdateHistory() error {
	return svc.DB.InsertHistoryFromCurrent()
}

// Cleanup remove old rows to save space
func (svc *SbsServiceImpl) Cleanup() error {
	return svc.DB.DeleteOldCurrent()
}
