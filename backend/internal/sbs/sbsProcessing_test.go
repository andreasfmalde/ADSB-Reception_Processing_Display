package sbs

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/models"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/mock"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	global.InitTestEnv()
	m.Run()
}

func TestProcessSbsStream(t *testing.T) {
	err := os.Chdir("../../")
	if err != nil {
		logger.Error.Fatalf("could not change working directory: %q", err)
	}

	mockLen3, err := os.ReadFile("./resources/mock/mockSbsDataLen3.txt")
	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	mockIncompleteData, err := os.ReadFile("./resources/mock/mockSbsIncompleteData.txt")
	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	mockParseError, err := os.ReadFile("./resources/mock/mockSbsParseError.txt")
	if err != nil {
		t.Errorf("error reading file: %q", err)
	}

	mockMalformedLines, err := os.ReadFile("./resources/mock/mockSbsMalformedDataLines.txt")
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
			expectedError:  errorMsg.ErrorCouldNotConnectToTcpStream,
			expectedLength: 3,
			mockResponse:   mockLen3,
		},
		{
			name:           "Incomplete Data Lines",
			expectedError:  errorMsg.ErrorCouldNotConnectToTcpStream,
			expectedLength: 0,
			mockResponse:   mockIncompleteData,
		},
		{
			name:           "Data Parsing Errors",
			expectedError:  errorMsg.ErrorCouldNotConnectToTcpStream,
			expectedLength: 0,
			mockResponse:   mockParseError,
		},
		{
			name:           "No Data Available",
			expectedError:  errorMsg.ErrorCouldNotConnectToTcpStream,
			expectedLength: 0,
			mockResponse:   []byte{},
		},
		{
			name:           "Malformed Data Lines",
			expectedError:  errorMsg.ErrorCouldNotConnectToTcpStream,
			expectedLength: 0,
			mockResponse:   mockMalformedLines,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := mock.InitStub(tt.mockResponse)
			addr, err := stub.StartServer()
			if err != nil {
				logger.Error.Fatalf("error starting stub server: %q", err)
				return
			}

			global.SbsSource = addr

			data, err := ProcessSBSstream()

			if (err != nil) != (tt.expectedError != "") || (err != nil && err.Error() != tt.expectedError) {
				t.Errorf("Test: %s Error = %s, expected %s", tt.name, err, tt.expectedError)
			}

			assert.Equal(t, tt.expectedLength, len(data))
		})
	}
}

func TestProcessSbsStream_ConnectionFailure(t *testing.T) {
	global.SbsSource = "unknown:5432"

	data, err := ProcessSBSstream()
	if err == nil && err.Error() != errorMsg.ErrorCouldNotConnectToTcpStream {
		t.Error("expected error")
	}

	assert.Equal(t, []models.AircraftCurrentModel{}, data)
}

func TestProcessSbsStream_Timeout(t *testing.T) {
	stub := mock.InitStub([]byte{})

	addr, err := stub.StartServer()
	if err != nil {
		logger.Error.Fatalf("error starting mock TCP server")
		return
	}
	global.SbsSource = addr

	stub.CloseConn()

	data, err := ProcessSBSstream()
	if err.Error() != errorMsg.ErrorCouldNotConnectToTcpStream {
		t.Errorf("error processing SBS data stream: %q", err)
	}

	assert.Nil(t, data)
}

func TestProcessSbsStream_ConnectionDrop(t *testing.T) {
	stub := mock.InitStub([]byte{})

	addr, err := stub.StartServer()
	if err != nil {
		logger.Error.Fatalf("error starting mock TCP server")
		return
	}
	global.SbsSource = addr

	stub.CloseListener()

	data, err := ProcessSBSstream()
	if err.Error() != errorMsg.ErrorCouldNotConnectToTcpStream {
		t.Errorf("error processing SBS data stream: %q", err)
	}

	assert.Nil(t, data)
}

func TestProcessSbsStream_ContinuousDataStream(t *testing.T) {

}

func generateMockAircraftResponses(n int) string {
	var builder strings.Builder

	const aircraftResponseTemplate = `
	MSG,1,0,0,E80451,0,2024/03/29,11:45:05.000,2024/03/29,11:45:05.000,TAM8112,,,,,,,,,,,
	MSG,3,0,0,E80451,0,2024/03/29,11:45:05.000,2024/03/29,11:45:05.000,,9725,,,19.329620,-99.196991,,,,,,
	MSG,4,0,0,E80451,0,2024/03/29,11:45:05.000,2024/03/29,11:45:05.000,,,184.317657,334.964325,,,-960,,,,,
	`

	for i := 0; i < n; i++ {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(aircraftResponseTemplate)
	}

	return builder.String()
}
