package global

// AircraftCurrentModel represents a row in current_time_aircraft
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

// AircraftHistoryModel represent a row in history_aircraft
type AircraftHistoryModel struct {
	Icao      string  `json:"icao"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

// GeoJson FeatureCollection for a Point type

type FeatureCollectionPoint struct {
	Type     string         `json:"type"`
	Features []FeaturePoint `json:"features"`
}

type FeaturePoint struct {
	Type       string                    `json:"type"`
	Properties AircraftCurrentProperties `json:"properties"`
	Geometry   GeometryPoint             `json:"geometry"`
}

type AircraftCurrentProperties struct {
	Icao         string `json:"icao"`
	Callsign     string `json:"callsign"`
	Altitude     int    `json:"altitude"`
	Speed        int    `json:"speed"`
	Track        int    `json:"track"`
	VerticalRate int    `json:"vspeed"`
	Timestamp    string `json:"timestamp"`
}

type GeometryPoint struct {
	Coordinates []float32 `json:"coordinates"`
	Type        string    `json:"type"`
}

// GeoJson FeatureCollection for a LineString type

type FeatureCollectionLineString struct {
	Type     string              `json:"type"`
	Features []FeatureLineString `json:"features"`
}

type FeatureLineString struct {
	Type       string                 `json:"type"`
	Properties AircraftHistProperties `json:"properties"`
	Geometry   GeometryLineString     `json:"geometry"`
}

type AircraftHistProperties struct {
	Icao string `json:"icao"`
}

type GeometryLineString struct {
	Coordinates [][]float32 `json:"coordinates"`
	Type        string      `json:"type"`
}
