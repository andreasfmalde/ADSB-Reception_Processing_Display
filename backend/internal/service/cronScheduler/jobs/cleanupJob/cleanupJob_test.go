package cleanupJob

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/mock"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)
	MaxDaysHistory := 5

	job := NewCleanupJob(mockDB, MaxDaysHistory)

	mockDB.EXPECT().DeleteOldHistory(MaxDaysHistory).Return(nil)

	job.Execute()
}

func TestCleanupJob_Execute_ErrorDeletingOldHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)
	MaxDaysHistory := 5

	job := NewCleanupJob(mockDB, MaxDaysHistory)

	var errorMessage = "mock error deleting old history data"

	mockDB.EXPECT().DeleteOldHistory(MaxDaysHistory).Return(errors.New(errorMessage))

	job.Execute()
}
