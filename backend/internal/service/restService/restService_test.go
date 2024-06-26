package restService

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

func Test_InitRestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	sbsSvc := InitRestService(mockDB)
	assert.NotNil(t, sbsSvc)
	assert.IsType(t, &RestImpl{}, sbsSvc)
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

	var errorMsg = "mockData error selecting table data"
	mockDB.EXPECT().SelectAllColumnsAircraftCurrent().Return(nil, errors.New(errorMsg))

	res, err := svc.GetCurrentAircraft()

	assert.NotNil(t, err)
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

	var errorMsg = "mockData error selecting table data"
	mockDB.EXPECT().SelectAllColumnHistoryByIcao("search").Return(nil, errors.New(errorMsg))

	res, err := svc.GetAircraftHistoryByIcao("search")

	assert.NotNil(t, err)
	assert.Equal(t, errorMsg, err.Error())
	assert.Nil(t, res)
}

func TestRestImpl_GetAircraftHistoryByIcaoFilterByTimestamp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &RestImpl{DB: mockDB}

	mockData := testUtility.CreateMockHistAircraft(10)
	search := mockData[0].Icao

	hour := 1

	mockDB.EXPECT().SelectAllColumnHistoryByIcaoFilterByTimestamp(search, hour).Return(mockData, nil)

	res, err := svc.GetAircraftHistoryByIcaoFilterByTimestamp(search, hour)

	assert.Nil(t, err)
	assert.Equal(t, mockData, res)
}

func TestRestImpl_GetAircraftHistoryByIcaoFilterByTimestamp_ErrorRetrievingDbData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockDatabase(ctrl)

	svc := &RestImpl{DB: mockDB}

	var errorMsg = "mockData error selecting table data"
	mockDB.EXPECT().SelectAllColumnHistoryByIcaoFilterByTimestamp("search", 1).Return(nil, errors.New(errorMsg))

	res, err := svc.GetAircraftHistoryByIcaoFilterByTimestamp("search", 1)

	assert.NotNil(t, err)
	assert.Equal(t, errorMsg, err.Error())
	assert.Nil(t, res)
}
