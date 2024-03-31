package convert

import (
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/geoJSON"
	"adsb-api/internal/global/models"
	"errors"
	"strconv"
	"strings"
)

func MakeTimeStamp(date string, time string) string {
	date = strings.Replace(date, "/", "-", -1)
	time = strings.TrimSuffix(time, ".000")
	return date + " " + time
}

func CurrentModelToGeoJson(aircraft []models.AircraftCurrentModel) (geoJSON.FeatureCollectionPoint, error) {
	var features []geoJSON.FeaturePoint
	for _, ac := range aircraft {
		var feature geoJSON.FeaturePoint
		feature.Type = "Feature"
		properties := geoJSON.AircraftCurrentProperties{
			Icao:         ac.Icao,
			Callsign:     ac.Callsign,
			Altitude:     ac.Altitude,
			Speed:        ac.Speed,
			Track:        ac.Track,
			VerticalRate: ac.VerticalRate,
			Timestamp:    ac.Timestamp,
		}
		feature.Properties = properties
		feature.Geometry.Type = "Point"
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, ac.Longitude, ac.Latitude)
		features = append(features, feature)
	}

	var featureCollection geoJSON.FeatureCollectionPoint
	featureCollection.Features = features
	featureCollection.Type = "FeatureCollection"
	return featureCollection, nil
}

func HistoryModelToGeoJson(aircraft []models.AircraftHistoryModel) (geoJSON.FeatureCollectionLineString, error) {
	if len(aircraft) < 2 {
		return geoJSON.FeatureCollectionLineString{}, errors.New(errorMsg.ErrorGeoJsonTooFewCoordinates)
	}

	var coordinates [][]float32
	for _, ac := range aircraft {
		point := []float32{ac.Latitude, ac.Longitude}
		coordinates = append(coordinates, point)
	}

	var features []geoJSON.FeatureLineString
	var feature geoJSON.FeatureLineString
	feature.Type = "Feature"
	feature.Properties.Icao = aircraft[0].Icao
	feature.Geometry.Coordinates = coordinates
	feature.Geometry.Type = "LineString"
	features = append(features, feature)

	var featureCollection geoJSON.FeatureCollectionLineString
	featureCollection.Features = features
	featureCollection.Type = "FeatureCollection"
	return featureCollection, nil
}

func SbsToAircraftCurrent(msg1 []string, msg3 []string, msg4 []string) (models.AircraftCurrentModel, error) {
	icao := msg1[4]
	date := msg1[8]
	hour := msg1[9]
	callsign := msg1[10]
	timestamp := MakeTimeStamp(date, hour)

	altitudeStr := msg3[11]
	latStr := msg3[14]
	longStr := msg3[15]

	altitude, err := strconv.Atoi(altitudeStr)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	lat, err := strconv.ParseFloat(latStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	long, err := strconv.ParseFloat(longStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}

	speedStr := msg4[12]
	trackStr := msg4[13]
	vspeedStr := msg4[16]

	speed, err := strconv.ParseFloat(speedStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	track, err := strconv.ParseFloat(trackStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	vspeed, err := strconv.Atoi(vspeedStr)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}

	return models.AircraftCurrentModel{
		Icao:         icao,
		Callsign:     callsign,
		Altitude:     altitude,
		Latitude:     float32(lat),
		Longitude:    float32(long),
		Speed:        int(speed),
		Track:        int(track),
		VerticalRate: vspeed,
		Timestamp:    timestamp,
	}, nil
}
