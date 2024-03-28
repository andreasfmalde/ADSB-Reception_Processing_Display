package service

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/mock"
	"adsb-api/internal/utility/testUtility"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	global.InitTestEnv()
	m.Run()
}

func Test_InitSbsService(t *testing.T) {
	sbsSvc, err := InitSbsService()
	if err != nil {
		t.Fatalf("error initiazling sbs eervice: %q", err)
	}

	assert.NotNil(t, sbsSvc)
}

func Test_InitSbsService_ErrorInitDB(t *testing.T) {
	global.DbUser = "user"
	sbsSvc, err := InitSbsService()
	if err == nil {
		t.Fatalf("error initiazling sbs eervice: %q", err)
	}

	assert.Nil(t, sbsSvc)
	global.DbUser = "test"
}

func TestSbsServiceImpl_CreateAdsbTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsServiceImpl{DB: mockDB}

	mockDB.EXPECT().BeginTx().Return(nil)
	mockDB.EXPECT().CreateAircraftCurrentTable().Return(nil)
	mockDB.EXPECT().CreateAircraftCurrentTimestampIndex().Return(nil)
	mockDB.EXPECT().Commit().Return(nil)
	mockDB.EXPECT().CreateAircraftHistoryTable().Return(nil)

	err := svc.CreateAdsbTables()

	assert.Nil(t, err)
}

func TestSbsServiceImpl_CreateAdsbTables_WithRollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsServiceImpl{DB: mockDB}

	var errorMsg = "mocking errorMsg creating table, should rollback transaction"

	mockDB.EXPECT().BeginTx().Return(nil)
	mockDB.EXPECT().CreateAircraftCurrentTable().Return(errors.New(errorMsg))
	mockDB.EXPECT().Rollback().Return(nil)

	err := svc.CreateAdsbTables()

	assert.Equal(t, errorMsg, err.Error())
}

func TestSbsServiceImpl_InsertNewSbsData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsServiceImpl{DB: mockDB}

	mockData := testUtility.CreateMockAircraft(100)

	mockDB.EXPECT().InsertHistoryFromCurrent().Return(nil)
	mockDB.EXPECT().BeginTx().Return(nil)
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

	svc := &SbsServiceImpl{DB: mockDB}

	mockData := testUtility.CreateMockAircraft(100)

	var errorMsg = "mocking errorMsg creating table, should rollback transaction"

	mockDB.EXPECT().InsertHistoryFromCurrent().Return(nil)
	mockDB.EXPECT().BeginTx().Return(nil)
	mockDB.EXPECT().DropAircraftCurrentTable().Return(errors.New(errorMsg))
	mockDB.EXPECT().Rollback().Return(nil)

	err := svc.InsertNewSbsData(mockData)

	assert.Equal(t, errorMsg, err.Error())
}

func TestSbsServiceImpl_Cleanup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsServiceImpl{DB: mockDB}

	mockDB.EXPECT().DeleteOldHistory(global.MaxDaysHistory).Return(nil)

	err := svc.Cleanup()

	assert.Nil(t, err)
}

func TestSbsServiceImpl_Cleanup_ErrorDeletingOldHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &SbsServiceImpl{DB: mockDB}

	var errorMsg = "mock error deleting old history data"

	mockDB.EXPECT().DeleteOldHistory(global.MaxDaysHistory).Return(errors.New(errorMsg))

	err := svc.Cleanup()

	assert.NotNil(t, err)
}
