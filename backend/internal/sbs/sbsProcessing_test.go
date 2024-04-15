package sbs

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/mock"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	m.Run()
}

func TestProcessSbsStream_WithMockResponse(t *testing.T) {
	err := os.Chdir("../../")
	if err != nil {
		t.Fatalf("could not change working directory: %q", err)
	}

	// 1 valid aircraft
	mockLen1, err := os.ReadFile("./resources/mockData/mockSbsDataLen1.txt")
	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	// missing MSG 4
	mockIncompleteData, err := os.ReadFile("./resources/mockData/mockSbsIncompleteData.txt")
	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	// speed value of MSG 2 is 'AAA'
	mockParseError, err := os.ReadFile("./resources/mockData/mockSbsParseError.txt")
	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	// data is valid but in incorrect order, MSG 4 before MSG 3
	mockMalformedLines, err := os.ReadFile("./resources/mockData/mockSbsMalformedDataLines.txt")
	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	mockFirstLineEmpty, err := os.ReadFile("./resources/mockData/mockSbsDataFirstLineEmpty.txt")

	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	tests := []struct {
		name           string
		mockResponse   []byte
		expectedError  string
		expectedLength int
	}{
		{
			name:           "Successful Data Retrieval",
			expectedLength: 1,
			mockResponse:   mockLen1,
		},
		{
			name:           "Successful Data Retrieval with big data",
			expectedLength: 1e5,
			mockResponse:   generateMockAircraftResponses(1e5),
		},
		{
			name:           "Incomplete Data Lines",
			expectedLength: 0,
			mockResponse:   mockIncompleteData,
		},
		{
			name:           "Data Parsing Errors",
			expectedLength: 0,
			mockResponse:   mockParseError,
		},
		{
			name:           "No Data Available",
			expectedLength: 0,
			mockResponse:   []byte{},
		},
		{
			name:           "Malformed Data Lines",
			expectedLength: 0,
			mockResponse:   mockMalformedLines,
		},
		{
			name:           "First line is empty",
			expectedLength: 0,
			mockResponse:   mockFirstLineEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := mock.InitStub(global.SbsSource, tt.mockResponse)
			err := stub.StartServer()
			if err != nil {
				t.Errorf("Test: %s Error = %s", tt.name, err)
				return
			}

			data, err := ProcessSbsStream(global.SbsSource, global.WaitingTime)
			if err != nil {
				t.Errorf("Test: %s Error = %s", tt.name, err)
			}

			assert.Equal(t, tt.expectedLength, len(data))

			if tt.expectedLength > 0 {
				for i, ac := range data {
					assert.NotEqualf(t, ac.Icao, "", "Test: %s Aircraft: %d Expected not nil: Icao", tt.name, i)
					assert.NotEqualf(t, ac.Callsign, "", "Test: %s Aircraft: %d Expected not nil: Callsign", tt.name, i)
					assert.NotEqualf(t, ac.Timestamp, "", "Test: %s Aircraft: %d Expected not nil: Timestamp", tt.name, i)
					assert.NotEqualf(t, ac.Altitude, 0, "Test: %s Aircraft: %d Expected not nil: Altitude", tt.name, i)
					assert.NotEqualf(t, ac.Latitude, 0, "Test: %s Aircraft: %d Expected not nil: Latitude", tt.name, i)
					assert.NotEqualf(t, ac.Longitude, 0, "Test: %s Aircraft: %d Expected not nil: Longitude", tt.name, i)
					assert.NotEqualf(t, ac.Track, 0, "Test: %s Aircraft: %d Expected not nil: Track", tt.name, i)
					assert.NotEqualf(t, ac.Speed, 0, "Test: %s Aircraft: %d Expected not nil: Speed", tt.name, i)
					assert.NotEqualf(t, ac.VerticalRate, 0, "Test: %s Aircraft: %d Expected not nil: VerticalRate", tt.name, i)
				}
			}
		})
	}
}

func TestProcessSbsStream_ConnectionFailure(t *testing.T) {
	global.SbsSource = "unknown:5432"

	data, err := ProcessSbsStream(global.SbsSource, global.WaitingTime)
	if err == nil {
		t.Error("expected error due to unknown host")
	}

	assert.Nil(t, data)

	// resets SbsSource const
	global.InitTestEnvironment()
}

func generateMockAircraftResponses(n int) []byte {
	var builder strings.Builder

	const responseLine1 = "MSG,1,0,0,E80451,0,2024/03/29,11:45:05.000,2024/03/29,11:45:05.000,TAM8112,,,,,,,,,,,"
	const responseLine2 = "MSG,3,0,0,E80451,0,2024/03/29,11:45:05.000,2024/03/29,11:45:05.000,,9725,,,19.329620,-99.196991,,,,,,"
	const responseLine3 = "MSG,4,0,0,E80451,0,2024/03/29,11:45:05.000,2024/03/29,11:45:05.000,,,184.317657,334.964325,,,-960,,,,,"

	for i := 0; i < n; i++ {
		builder.WriteString(responseLine1 + "\n")
		builder.WriteString(responseLine2 + "\n")
		builder.WriteString(responseLine3 + "\n")
	}

	return []byte(builder.String())
}
