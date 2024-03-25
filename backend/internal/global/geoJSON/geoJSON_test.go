package geoJSON

import (
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/testUtility"
	"github.com/xeipuuv/gojsonschema"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

var geoJsonOverallSchema string

func TestMain(m *testing.M) {
	logger.InitLogger()
	err := os.Chdir("../../../")
	if err != nil {
		logger.Error.Fatalf("could not change working directory: %q", err)
	}

	absPath, err := filepath.Abs("./resources/schemas/geoJson.json")
	if err != nil {
		logger.Error.Fatalf("could not get absolute path: %q", err)
	}

	absPath = filepath.ToSlash(absPath)

	u := url.URL{}
	u.Scheme = "file"
	u.Path = absPath
	absURL := u.String()
	geoJsonOverallSchema = absURL

	m.Run()

}

func TestConvertCurrentModelToGeoJson(t *testing.T) {
	schemaLoader := gojsonschema.NewReferenceLoader(geoJsonOverallSchema)

	var mockData = testUtility.CreateMockAircraft(1)
	geoJson, err := ConvertCurrentModelToGeoJson(mockData)
	if err != nil {
		t.Errorf("error converting model data to GeoJSON: %q", err)
	}

	documentLoader := gojsonschema.NewGoLoader(geoJson)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Errorf("Error validating response body: %s", err.Error())
	}

	if !result.Valid() {
		t.Errorf("Response body does not follow the GeoJSON standard")
		for _, desc := range result.Errors() {
			t.Logf("- %s", desc)
		}
	}
}

func TestConvertHistoryModelToGeoJson(t *testing.T) {
	schemaLoader := gojsonschema.NewReferenceLoader(geoJsonOverallSchema)

	var mockData = testUtility.CreateMockHistAircraft(2)
	geoJson, err := ConvertHistoryModelToGeoJson(mockData)
	if err != nil {
		t.Errorf("error converting model data to GeoJSON: %q", err)
	}

	documentLoader := gojsonschema.NewGoLoader(geoJson)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Errorf("Error validating response body: %s", err.Error())
	}

	if !result.Valid() {
		t.Errorf("Response body does not follow the GeoJSON standard")
		for _, desc := range result.Errors() {
			t.Logf("- %s", desc)
		}
	}
}

func TestConvertHistoryModelToGeoJson_TooFewCoordinates(t *testing.T) {
	var mockData = testUtility.CreateMockHistAircraft(1)
	_, err := ConvertHistoryModelToGeoJson(mockData)
	if err == nil {
		t.Errorf("expected error: %s", global.ErrorGeoJsonTooFewCoordinates)
	}
}
