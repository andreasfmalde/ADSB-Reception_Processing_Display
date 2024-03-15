package currentAircraftHandler

import (
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var currentEndpoint *httptest.Server

func TestMain(m *testing.M) {
	logger.InitLogger()
	dbConn := utility.InitTestDb()
	currentEndpoint = httptest.NewServer(CurrentAircraftHandler(dbConn))
	defer currentEndpoint.Close()

	err := os.Chdir(filepath.Dir("../../go.mod"))
	if err != nil {
		logger.Error.Fatalf("unable to retrieve file.")
	}

	m.Run()
}

func TestValidRequests(t *testing.T) {
	output, err := os.ReadFile("./../resources/mockAircraft/aircraft.json")
	if err != nil {
		log.Printf("error while reading file: %v\n", err)
		return
	}

	var aircraft global.GeoJsonFeatureCollection
	err = json.Unmarshal(output, &aircraft)
	if err != nil {
		logger.Error.Fatalf(err.Error())
	}

	tests := []struct {
		name, url, method     string
		expectedStatusCode    int
		expectedResponse      global.GeoJsonFeatureCollection
		expectedResponseArray []global.GeoJsonFeatureCollection
	}{
		{
			name:               "Get request without parameters",
			url:                currentEndpoint.URL + global.CurrentAircraftPath,
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Errorf("Test: %s. Error creating request with method %s and url %s: %s", tt.name, tt.method, tt.url, err.Error())
			}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Test: %s. Error executing %s request: %s", tt.name, tt.method, err.Error())
			}
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}

}
