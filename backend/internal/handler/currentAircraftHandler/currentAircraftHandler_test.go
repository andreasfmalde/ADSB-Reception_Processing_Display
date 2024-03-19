package currentAircraftHandler

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var currentEndpoint *httptest.Server
var dbConn db.AdsbDB

func TestMain(m *testing.M) {
	logger.InitLogger()
	global.InitEnv()
	adsbDB, err := db.InitDB()
	if err != nil {
		logger.Error.Fatalf("Error opening database: %q", err)
	}
	dbConn = *adsbDB

	adsbSvc := db.AdsbService{DB: adsbDB}

	currentEndpoint = httptest.NewServer(CurrentAircraftHandler(&adsbSvc))
	defer currentEndpoint.Close()
	m.Run()
}

func TestInvalidRequests(t *testing.T) {
	tests := []struct {
		name, url, method string
		statusCode        int
		setup             func()
		teardown          func()
	}{
		{
			name:       "Post request",
			url:        currentEndpoint.URL + global.CurrentAircraftPath,
			method:     http.MethodPost,
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "Delete request",
			url:        currentEndpoint.URL + global.CurrentAircraftPath,
			method:     http.MethodDelete,
			statusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if tt.teardown != nil {
				defer tt.teardown()
			}

			client := &http.Client{}
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Errorf("Test: %s. Error creating request with populateMethod %s and url %s: %s", tt.name, tt.method, tt.url, err.Error())
			}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Test: %s. Error executing %s request: %s", tt.name, tt.method, err.Error())
			}
			assert.Equal(t, tt.statusCode, res.StatusCode)
		})
	}
}

func TestValidRequests(t *testing.T) {

	tests := []struct {
		name, url, httpMethod string
		statusCode, length    int
		setup                 func()
		teardown              func()
	}{
		{
			name:       "Get request without parameters",
			url:        currentEndpoint.URL + global.CurrentAircraftPath,
			httpMethod: http.MethodGet,
			statusCode: http.StatusOK,
			length:     10,
			setup: func() {

			},
			teardown: func() {

			},
		},
		{
			name:       "Get request with empty database",
			url:        currentEndpoint.URL + global.CurrentAircraftPath,
			httpMethod: http.MethodGet,
			statusCode: http.StatusNotFound,
			length:     0,
			setup: func() {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if tt.teardown != nil {
				defer tt.teardown()
			}

			client := &http.Client{}
			req, err := http.NewRequest(tt.httpMethod, tt.url, nil)
			if err != nil {
				t.Errorf("Test: %s. Error creating request with populateMethod %s and url %s: %s", tt.name, tt.httpMethod, tt.url, err.Error())
			}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Test: %s. Error executing %s request: %s", tt.name, tt.httpMethod, err.Error())
			}

			assert.Equal(t, tt.statusCode, res.StatusCode)

			if tt.length > 0 {
				var actual global.GeoJsonFeatureCollection
				err = json.NewDecoder(res.Body).Decode(&actual)
				if err != nil {
					t.Errorf("Error decoding response body: %s", err.Error())
				}

				assert.Equal(t, tt.length, len(actual.Features))
			}
		})
	}
}
