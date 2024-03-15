package global

/*
Struct used to represent an aircraft record
*/
type Aircraft struct {
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

type GeoJsonFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJsonFeature `json:"features"`
}

type GeoJsonFeature struct {
	Type       string             `json:"type"`
	Properties AircraftProperties `json:"properties"`
	Geometry   struct {
		Coordinates []float32 `json:"coordinates"`
		Type        string    `json:"type"`
	} `json:"geometry"`
}

type AircraftProperties struct {
	Icao         string `json:"icao"`
	Callsign     string `json:"callsign"`
	Altitude     int    `json:"altitude"`
	Speed        int    `json:"speed"`
	Track        int    `json:"track"`
	VerticalRate int    `json:"vspeed"`
	Timestamp    string `json:"timestamp"`
}
