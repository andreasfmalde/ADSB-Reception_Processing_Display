package db

import (
	"adsb-api/internal/global"
	"testing"
)

func TestMain(m *testing.M) {
	global.InitDevEnv()
	m.Run()
}

func TestInitDatabase(t *testing.T) {
	_, err := Init()
	if err != nil {
		t.Errorf("Database connection failed: %q", err)
	}
}

func TestBulkInsertCurrentAircraftTable(t *testing.T) {

}
