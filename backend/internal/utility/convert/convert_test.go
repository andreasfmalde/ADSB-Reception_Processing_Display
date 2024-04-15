package convert

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/utility/testUtility"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

var schemaLoader gojsonschema.JSONLoader

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	err := os.Chdir("../../../")
	if err != nil {
		log.Fatalf("could not change working directory: %q", err)
	}

	absPath, err := filepath.Abs("./resources/schemas/geoJson.json")
	if err != nil {
		log.Fatalf("could not get absolute path: %q", err)
	}

	absPath = filepath.ToSlash(absPath)

	u := url.URL{}
	u.Scheme = "file"
	u.Path = absPath
	absURL := u.String()
	schemaLoader = gojsonschema.NewReferenceLoader(absURL)

	m.Run()
}

func TestConvertCurrentModelToGeoJson(t *testing.T) {
	var mockData = testUtility.CreateMockAircraft(1)
	geoJson, err := CurrentModelToGeoJson(mockData)
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
	var mockData = testUtility.CreateMockHistAircraft(2)
	geoJson, err := HistoryModelToGeoJson(mockData)
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
	_, err := HistoryModelToGeoJson(mockData)
	if err == nil {
		t.Errorf("expected error: %s", errorMsg.ErrorGeoJsonTooFewCoordinates)
	}
}
