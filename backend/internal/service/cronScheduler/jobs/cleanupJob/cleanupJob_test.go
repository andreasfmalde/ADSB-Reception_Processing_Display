package cleanupJob

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/utility/mock"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	m.Run()
}

func TestNewCleanupJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	job := NewCleanupJob(mockDB, 5)

	assert.NotNil(t, job, job.db, job.MaxDaysHistory)
}

func TestCleanupJob_Execute(t *testing.T) {
	var logBuffer bytes.Buffer
	log.Logger = zerolog.New(&logBuffer)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)
	MaxDaysHistory := 5

	job := NewCleanupJob(mockDB, MaxDaysHistory)

	mockDB.EXPECT().DeleteOldHistory(MaxDaysHistory).Return(nil)

	job.Execute()

	logOutput := logBuffer.String()

	assert.Contains(t, logOutput, errorMsg.InfoOldHistoryDataDeleted)
	log.Logger = zerolog.New(os.Stderr)
}

func TestCleanupJob_Execute_ErrorDeletingOldHistory(t *testing.T) {
	var logBuffer bytes.Buffer
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: &logBuffer, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Logger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)
	MaxDaysHistory := 5

	job := NewCleanupJob(mockDB, MaxDaysHistory)

	var errorMessage = "mockData error deleting old history data"

	mockDB.EXPECT().DeleteOldHistory(MaxDaysHistory).Return(errors.New(errorMessage))

	job.Execute()

	logOutput := logBuffer.String()

	assert.Contains(t, logOutput, fmt.Errorf(errorMsg.ErrorDeletingOldHistory+": %q", errorMessage).Error())
	log.Logger = zerolog.New(os.Stderr)
}
