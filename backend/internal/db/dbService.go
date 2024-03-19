package db

import (
	"adsb-api/internal/global"
)

type AdsbServiceInterface interface {
	InitSvc() (*AdsbService, error)
	Close() error
	CreateCurrentTimeAircraftTable() error
	BulkInsertCurrentTimeAircraftTable(aircraft []global.Aircraft) error
	DeleteOldCurrentAircraft() error
	GetAllCurrentAircraft() (global.GeoJsonFeatureCollection, error)
}

type AdsbService struct {
	DB *AdsbRepository
}

func (service *AdsbService) InitSvc() (*AdsbService, error) {
	repo, err := InitDB()
	if err != nil {
		return nil, err
	}
	return &AdsbService{DB: repo}, nil
}

func (service *AdsbService) Close() error {
	return service.DB.Close()
}

func (service *AdsbService) CreateCurrentTimeAircraftTable() error {
	return service.DB.CreateCurrentTimeAircraftTable()
}

func (service *AdsbService) BulkInsertCurrentTimeAircraftTable(aircraft []global.Aircraft) error {
	return service.DB.BulkInsertCurrentTimeAircraftTable(aircraft)
}

func (service *AdsbService) DeleteOldCurrentAircraft() error {
	return service.DB.DeleteOldCurrentAircraft()
}

func (service *AdsbService) GetAllCurrentAircraft() (global.GeoJsonFeatureCollection, error) {
	return service.DB.GetAllCurrentAircraft()
}
