package global

import (
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/mock"
	"github.com/joho/godotenv"
	"os"
)

func InitEnvironment() {
	_ = godotenv.Load("../.env")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
}

func InitTestEnv() {
	logger.InitLogger()
	DbUser = "test"
	DbPassword = "test"
	Dbname = "adsb_test_db"

	SbsSource = "localhost:9999"
}

func StartStubServer() {
	mockData, err := os.ReadFile("./resources/mock/mockSbsDataLen5.txt")
	if err != nil {
		logger.Error.Printf("error reading file: %q", err)
	}

	stub := mock.InitStub(SbsSource, mockData)
	err = stub.StartServer()
	if err != nil {
		logger.Error.Fatalf("error starting stub server: %q", err)
		return
	}

}
