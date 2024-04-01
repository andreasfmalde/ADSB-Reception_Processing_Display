package service

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/models"
)

type RestService interface {
	GetCurrentAircraft() ([]models.AircraftCurrentModel, error)
	GetHistoryByIcao(search string) ([]models.AircraftHistoryModel, error)
}

type RestServiceImpl struct {
	DB db.Database
}

// InitRestService initializes RestServiceImpl struct and database connection
func InitRestService() (*RestServiceImpl, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	return &RestServiceImpl{DB: dbConn}, nil
}

// GetCurrentAircraft retrieves a list of all aircraft that are considered 'current'
// (i.e., aircraft that are currently in the air).
func (svc *RestServiceImpl) GetCurrentAircraft() ([]models.AircraftCurrentModel, error) {
	return svc.DB.GetCurrentAircraft()
}

// GetHistoryByIcao retrieves aircraft history from given icao.
func (svc *RestServiceImpl) GetHistoryByIcao(icao string) ([]models.AircraftHistoryModel, error) {
	return svc.DB.GetHistoryByIcao(icao)
}
