package restService

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/models"
)

type RestService interface {
	GetCurrentAircraft() ([]models.AircraftCurrentModel, error)
	GetAircraftHistoryByIcao(search string) ([]models.AircraftHistoryModel, error)
}

type RestImpl struct {
	DB db.PgDatabase
}

// InitRestService initializes RestImpl struct and database connection
func InitRestService() (*RestImpl, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return nil, err
	}
	return &RestImpl{DB: dbConn}, nil
}

// GetCurrentAircraft retrieves a list of all aircraft that are considered 'current'
// (i.e., aircraft that are currently in the air).
func (svc *RestImpl) GetCurrentAircraft() ([]models.AircraftCurrentModel, error) {
	return svc.DB.SelectAllColumnsAircraftCurrent()
}

// GetAircraftHistoryByIcao retrieves aircraft history from given icao.
func (svc *RestImpl) GetAircraftHistoryByIcao(icao string) ([]models.AircraftHistoryModel, error) {
	return svc.DB.SelectAllColumnHistoryByIcao(icao)
}
