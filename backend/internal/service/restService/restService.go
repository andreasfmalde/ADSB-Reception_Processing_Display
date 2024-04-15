package restService

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global/models"
)

// RestService is an interface representing a RESTful service for retrieving database data through the repository in
// internal/db/database.go
type RestService interface {
	GetCurrentAircraft() ([]models.AircraftCurrentModel, error)
	GetAircraftHistoryByIcao(search string) ([]models.AircraftHistoryModel, error)
	GetAircraftHistoryByIcaoFilterByTimestamp(search string, hour int) ([]models.AircraftHistoryModel, error)
}

type RestImpl struct {
	DB db.Database
}

// InitRestService initializes RestImpl struct and database connection.
func InitRestService(db db.Database) *RestImpl {
	return &RestImpl{DB: db}
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

// GetAircraftHistoryByIcaoFilterByTimestamp retrieves aircraft by ICAO code and limits the results by only retrieving
// data newer than given hour parameter.
func (svc *RestImpl) GetAircraftHistoryByIcaoFilterByTimestamp(search string, hour int) ([]models.AircraftHistoryModel, error) {
	return svc.DB.SelectAllColumnHistoryByIcaoFilterByTimestamp(search, hour)
}
