package testUtility

import (
	"adsb-api/internal/global"
	"strconv"
	"time"
)

func CreateMockAircraft(n int) []global.Aircraft {
	var aircraft []global.Aircraft

	for i := 0; i < n; i++ {
		ac := global.Aircraft{
			Icao:         strconv.Itoa(i),
			Callsign:     strconv.Itoa(i),
			Altitude:     i,
			Latitude:     float32(i),
			Longitude:    float32(i),
			Speed:        i,
			Track:        i,
			VerticalRate: i,
			Timestamp:    time.Now().Format(time.DateTime),
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}

func CreateMockAircraftWithTimestamp(icao string, timestamp string) global.Aircraft {
	return global.Aircraft{
		Icao:         icao,
		Callsign:     "TEST",
		Altitude:     10000,
		Latitude:     51.5074,
		Longitude:    0.1278,
		Speed:        450,
		Track:        180,
		VerticalRate: 0,
		Timestamp:    timestamp,
	}
}

func CreateMockFeatureCollection(n int) global.GeoJsonFeatureCollection {
	featureCollection := global.GeoJsonFeatureCollection{}
	featureCollection.Type = "FeatureCollection"

	for i := 0; i < n; i++ {
		var lat float32 = 51.5074
		var long float32 = 51.5074

		ac := global.GeoJsonFeature{
			Type: "Feature",
			Properties: global.AircraftProperties{
				Icao:         "TEST",
				Callsign:     "TEST",
				Altitude:     0,
				Speed:        0,
				Track:        0,
				VerticalRate: 0,
				Timestamp:    "",
			},
			Geometry: struct {
				Coordinates []float32 `json:"coordinates"`
				Type        string    `json:"type"`
			}{},
		}

		feature := global.GeoJsonFeature{}
		feature.Properties = ac.Properties
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, lat, long)
		feature.Geometry.Type = "Point"

		featureCollection.Features = append(featureCollection.Features, feature)
	}

	return featureCollection
}
