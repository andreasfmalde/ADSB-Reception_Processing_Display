package testUtility

import (
	"adsb-api/internal/global"
	"strconv"
	"time"
)

func CreateMockAircraft(n int) []global.AircraftCurrentModel {
	var aircraft []global.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := global.AircraftCurrentModel{
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

func CreateMockAircraftWithTimestamp(icao string, timestamp string) global.AircraftCurrentModel {
	return global.AircraftCurrentModel{
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

func CreateMockAircraftWithIcao(n int, icao string) []global.AircraftCurrentModel {
	var aircraft []global.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := global.AircraftCurrentModel{
			Icao:         icao,
			Callsign:     strconv.Itoa(i),
			Altitude:     i,
			Latitude:     float32(i),
			Longitude:    float32(i),
			Speed:        i,
			Track:        i,
			VerticalRate: i,
			Timestamp:    time.Now().Add(time.Duration(i) * time.Second).Format(time.DateTime),
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}

func CreateMockFeatureCollectionPoint(n int) global.FeatureCollectionPoint {
	featureCollection := global.FeatureCollectionPoint{}
	featureCollection.Type = "FeatureCollection"

	for i := 0; i < n; i++ {
		var lat float32 = 51.5074
		var long float32 = 51.5074

		ac := global.FeaturePoint{
			Type: "Feature",
			Properties: global.AircraftCurrentProperties{
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

		feature := global.FeaturePoint{}
		feature.Properties = ac.Properties
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, lat, long)
		feature.Geometry.Type = "Point"

		featureCollection.Features = append(featureCollection.Features, feature)
	}

	return featureCollection
}

func CreateMockFeatureCollectionLineString(n int) global.FeatureCollectionLineString {
	var coordinates [][]float32

	for i := 0; i < n; i++ {
		coordinates = append(coordinates, []float32{float32(i), float32(-i)})
	}

	feature := global.FeatureLineString{}
	feature.Type = "Feature"
	feature.Properties.Icao = "TEST"
	feature.Geometry.Type = "LineString"
	feature.Geometry.Coordinates = coordinates

	featureCollection := global.FeatureCollectionLineString{}
	featureCollection.Type = "FeatureCollection"
	featureCollection.Features = append(featureCollection.Features, feature)

	return featureCollection
}
