package models

// AircraftHistoryModel represent a row in history_aircraft
type AircraftHistoryModel struct {
	Icao      string  `json:"icao"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}
