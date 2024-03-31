package geoJSON

// GeoJson FeatureCollection for a Point type

type FeatureCollectionPoint struct {
	Type     string         `json:"type"`
	Features []FeaturePoint `json:"features"`
}

type FeaturePoint struct {
	Type       string                    `json:"type"`
	Geometry   geometryPoint             `json:"geometry"`
	Properties AircraftCurrentProperties `json:"properties"`
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

type geometryPoint struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

// GeoJson FeatureCollection for a LineString type

type FeatureCollectionLineString struct {
	Type     string              `json:"type"`
	Features []FeatureLineString `json:"features"`
}

type FeatureLineString struct {
	Type       string                 `json:"type"`
	Properties aircraftHistProperties `json:"properties"`
	Geometry   geometryLineString     `json:"geometry"`
}

type aircraftHistProperties struct {
	Icao string `json:"icao"`
}

type geometryLineString struct {
	Coordinates [][]float32 `json:"coordinates"`
	Type        string      `json:"type"`
}
