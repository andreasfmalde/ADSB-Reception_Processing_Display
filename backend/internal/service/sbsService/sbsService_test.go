package sbsService

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/mock"
	"adsb-api/internal/utility/testUtility"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	m.Run()
}

func Test_InitSbsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)
	mockCron := mock.NewMockScheduler(ctrl)

	sbsSvc := InitSbsService(mockDB, mockCron)

	assert.NotNil(t, sbsSvc, sbsSvc.DB, sbsSvc.CronScheduler)
}

func TestSbsServiceImpl_CreateAdsbTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsImpl{DB: mockDB}

	mockDB.EXPECT().CreateAircraftCurrentTable().Return(nil)
	mockDB.EXPECT().Begin().Return(nil)
	mockDB.EXPECT().CreateAircraftHistoryTable().Return(nil)
	mockDB.EXPECT().CreateAircraftHistoryTimestampIndex().Return(nil)
	mockDB.EXPECT().Commit().Return(nil)
	err := svc.CreateAdsbTables()

	assert.Nil(t, err)
}

func TestSbsServiceImpl_CreateAdsbTables_WithRollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsImpl{DB: mockDB}

	var errorMsg = "mocking errorMsg creating table, should rollback transaction"

	mockDB.EXPECT().CreateAircraftCurrentTable().Return(nil)
	mockDB.EXPECT().Begin().Return(nil)
	mockDB.EXPECT().CreateAircraftHistoryTable().Return(nil)
	mockDB.EXPECT().CreateAircraftHistoryTimestampIndex().Return(errors.New(errorMsg))
	mockDB.EXPECT().Rollback().Return(nil)

	err := svc.CreateAdsbTables()

	assert.Equal(t, errorMsg, err.Error())
}

func TestSbsServiceImpl_InsertNewSbsData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsImpl{DB: mockDB}

	mockData := testUtility.CreateMockAircraft(100)

	mockDB.EXPECT().InsertHistoryFromCurrent().Return(nil)
	mockDB.EXPECT().Begin().Return(nil)
	mockDB.EXPECT().DropAircraftCurrentTable().Return(nil)
	mockDB.EXPECT().CreateAircraftCurrentTable().Return(nil)
	mockDB.EXPECT().BulkInsertAircraftCurrent(mockData).Return(nil)
	mockDB.EXPECT().Commit().Return(nil)

	err := svc.InsertNewSbsData(mockData)

	assert.Nil(t, err)
}

func TestSbsServiceImpl_InsertNewSbsData_WithRollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsImpl{DB: mockDB}

	mockData := testUtility.CreateMockAircraft(100)

	var errorMsg = "mocking errorMsg creating table, should rollback transaction"

	mockDB.EXPECT().InsertHistoryFromCurrent().Return(nil)
	mockDB.EXPECT().Begin().Return(nil)
	mockDB.EXPECT().DropAircraftCurrentTable().Return(errors.New(errorMsg))
	mockDB.EXPECT().Rollback().Return(nil)

	err := svc.InsertNewSbsData(mockData)

	assert.NotNil(t, err)
	assert.Equal(t, errorMsg, err.Error())
}

func TestSbsImpl_ScheduleCleanUpJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)
	mockCron := mock.NewMockScheduler(ctrl)
	svc := &SbsImpl{DB: mockDB, CronScheduler: mockCron}

	schedule := "* * * * *"
	MaxDaysHistory := 5

	mockCron.EXPECT().ScheduleJob(schedule, gomock.Any()).Return(nil)

	err := svc.ScheduleCleanUpJob(schedule, MaxDaysHistory)

	assert.Nil(t, err)
}

func TestSbsImpl_ScheduleCleanUpJob_ErrorScheduleJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)
	mockCron := mock.NewMockScheduler(ctrl)
	svc := &SbsImpl{DB: mockDB, CronScheduler: mockCron}

	schedule := "* * * * *"
	MaxDaysHistory := 5

	errorMessage := "mockData error simulating error scheduling job"

	mockCron.EXPECT().ScheduleJob(schedule, gomock.Any()).Return(errors.New(errorMessage))

	err := svc.ScheduleCleanUpJob(schedule, MaxDaysHistory)

	assert.NotNil(t, err)
	assert.Equal(t, errorMessage, err.Error())
}
