// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\service\sbsService.go

// Package service is a generated GoMock package.
package mock

import (
	models "adsb-api/internal/global/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSbsService is a mock of SbsService interface.
type MockSbsService struct {
	ctrl     *gomock.Controller
	recorder *MockSbsServiceMockRecorder
}

// MockSbsServiceMockRecorder is the mock recorder for MockSbsService.
type MockSbsServiceMockRecorder struct {
	mock *MockSbsService
}

// NewMockSbsService creates a new mock instance.
func NewMockSbsService(ctrl *gomock.Controller) *MockSbsService {
	mock := &MockSbsService{ctrl: ctrl}
	mock.recorder = &MockSbsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSbsService) EXPECT() *MockSbsServiceMockRecorder {
	return m.recorder
}

// Cleanup mocks base method.
func (m *MockSbsService) Cleanup() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cleanup")
	ret0, _ := ret[0].(error)
	return ret0
}

// Cleanup indicates an expected call of Cleanup.
func (mr *MockSbsServiceMockRecorder) Cleanup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockSbsService)(nil).Cleanup))
}

// CreateAdsbTables mocks base method.
func (m *MockSbsService) CreateAdsbTables() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdsbTables")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAdsbTables indicates an expected call of CreateAdsbTables.
func (mr *MockSbsServiceMockRecorder) CreateAdsbTables() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdsbTables", reflect.TypeOf((*MockSbsService)(nil).CreateAdsbTables))
}

// InsertNewSbsData mocks base method.
func (m *MockSbsService) InsertNewSbsData(aircraft []models.AircraftCurrentModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewSbsData", aircraft)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewSbsData indicates an expected call of InsertNewSbsData.
func (mr *MockSbsServiceMockRecorder) InsertNewSbsData(aircraft interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewSbsData", reflect.TypeOf((*MockSbsService)(nil).InsertNewSbsData), aircraft)
}