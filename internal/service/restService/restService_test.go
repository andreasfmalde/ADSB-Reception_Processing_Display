package restService

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
	global.InitTestEnvironment()
	m.Run()
}

func Test_InitRestService(t *testing.T) {
	sbsSvc, err := InitRestService()
	if err != nil {
		t.Fatalf("error initiazling sbs eervice: %q", err)
	}

	assert.NotNil(t, sbsSvc)
}

func Test_InitRestService_ErrorInitDB(t *testing.T) {
	global.DbUser = "user"
	sbsSvc, err := InitRestService()
	if err == nil {
		t.Fatalf("error initiazling sbs eervice: %q", err)
	}

	assert.Nil(t, sbsSvc)
	global.DbUser = "test"
}

func TestRestServiceImpl_GetCurrentAircraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &RestImpl{DB: mockDB}

	mockData := testUtility.CreateMockAircraft(100)
	mockDB.EXPECT().SelectAllColumnsAircraftCurrent().Return(mockData, nil)

	res, err := svc.GetCurrentAircraft()

	assert.Nil(t, err)
	assert.Equal(t, mockData, res)

}

func TestRestServiceImpl_GetCurrentAircraft_ErrorRetrievingDbData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &RestImpl{DB: mockDB}

	var errorMsg = "mock error selecting table data"
	mockDB.EXPECT().SelectAllColumnsAircraftCurrent().Return(nil, errors.New(errorMsg))

	res, err := svc.GetCurrentAircraft()

	assert.Equal(t, errorMsg, err.Error())
	assert.Nil(t, res)
}

func TestRestServiceImpl_GetAircraftHistoryByIcao(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &RestImpl{DB: mockDB}

	mockData := testUtility.CreateMockHistAircraft(10)
	var search = mockData[0].Icao

	mockDB.EXPECT().SelectAllColumnHistoryByIcao(search).Return(mockData, nil)

	res, err := svc.GetAircraftHistoryByIcao(search)

	assert.Nil(t, err)
	assert.Equal(t, mockData, res)
}

func TestRestServiceImpl_GetAircraftHistoryByIcao_ErrorRetrievingDbData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &RestImpl{DB: mockDB}

	var errorMsg = "mock error selecting table data"
	mockDB.EXPECT().SelectAllColumnHistoryByIcao("search").Return(nil, errors.New(errorMsg))

	res, err := svc.GetAircraftHistoryByIcao("search")

	assert.Equal(t, errorMsg, err.Error())
	assert.Nil(t, res)
}
