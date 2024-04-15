package models

// AircraftHistoryModel represent a row in aircraft_history
type AircraftHistoryModel struct {
	Icao      string  `json:"icao"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

// AircraftCurrentModel represents a row in aircraft_current
type AircraftCurrentModel struct {
	Icao         string  `json:"icao"`
	Callsign     string  `json:"callsign"`
	Altitude     int     `json:"altitude"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	Speed        int     `json:"speed"`
	Track        int     `json:"track"`
	VerticalRate int     `json:"vspeed"`
	Timestamp    string  `json:"timestamp"`
}
