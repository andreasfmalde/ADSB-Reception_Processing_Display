package db

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/models"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/testUtility"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	global.InitTestEnv()
	m.Run()
}

func setupTestDB() *AdsbDB {
	db, err := InitDB()
	if err != nil {
		logger.Error.Fatalf("Failed to initialize service: %v", err)
	}

	err = db.createAircraftCurrentTable()
	if err != nil {
		logger.Error.Fatalf("error creating current_time_aircraft table: %q", err)
	}

	err = db.createAircraftHistory()
	if err != nil {
		logger.Error.Fatalf("error creating history_aircraft table: %q", err)
	}

	return db
}

func teardownTestDB(db *AdsbDB) {
	dropAircraftCurrent(db)
	dropAircraftHistory(db)

	err := db.Close()
	if err != nil {
		logger.Error.Fatalf("error closing database: %q", err)
	}
}

func dropAircraftCurrent(db *AdsbDB) {
	_, err := db.Conn.Exec("DROP TABLE IF EXISTS aircraft_current CASCADE")
	if err != nil {
		logger.Error.Fatalf("error droppint current_time_aircraft: %q", err.Error())
	}
}

func dropAircraftHistory(db *AdsbDB) {
	_, err := db.Conn.Exec("DROP TABLE IF EXISTS aircraft_history CASCADE")
	if err != nil {
		logger.Error.Fatalf("error droppint current_time_aircraft: %q", err.Error())
	}
}

func TestInitCloseDB(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Errorf("Database connection failed: %q", err)
	}
	defer func(db *AdsbDB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	err = db.Close()
	if err != nil {
		t.Errorf("Database connection failed: %q", err)
	}
}

func TestAdsbDB_CreateAdsbTables(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	dropAircraftCurrent(db)
	dropAircraftHistory(db)

	err := db.CreateAdsbTables()
	if err != nil {
		t.Errorf("creating ADS-B tables failed: %q", err)
	}

	expectedCurrentTimeAircraftColumns := map[string]string{
		"icao":      "character varying(6)",
		"callsign":  "character varying(10)",
		"altitude":  "integer",
		"lat":       "numeric",
		"long":      "numeric",
		"speed":     "integer",
		"track":     "integer",
		"vspeed":    "integer",
		"timestamp": "timestamp without time zone",
	}

	expectedHistoryAircraftColumns := map[string]string{
		"icao":      "character varying(6)",
		"lat":       "numeric",
		"long":      "numeric",
		"timestamp": "timestamp without time zone",
	}

	checkTableColumns(t, db, "aircraft_current", expectedCurrentTimeAircraftColumns)
	checkTableColumns(t, db, "aircraft_history", expectedHistoryAircraftColumns)
}

func checkTableColumns(t *testing.T, db *AdsbDB, tableName string, expectedColumns map[string]string) {
	rows, err := db.Conn.Query(`SELECT column_name, data_type, character_maximum_length FROM information_schema.columns WHERE table_name = $1`, tableName)
	if err != nil {
		t.Fatalf("error executing column info query: %q", err.Error())
	}
	defer rows.Close()

	actualColumns := make(map[string]string)
	for rows.Next() {
		var columnName, dataType string
		var maxLength *int
		if err := rows.Scan(&columnName, &dataType, &maxLength); err != nil {
			t.Fatalf("error scanning column info: %q", err.Error())
		}
		if maxLength != nil {
			dataType = fmt.Sprintf("%s(%d)", dataType, *maxLength)
		}
		actualColumns[columnName] = dataType
	}

	if !reflect.DeepEqual(expectedColumns, actualColumns) {
		t.Fatalf("columns for table %s do not match expected. \n Expected: %v \n Got     : %v", tableName, expectedColumns, actualColumns)
	}
}

func TestAdsbDB_BulkInsertAircraftCurrent(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	var nAircraft = 100

	aircraft := testUtility.CreateMockAircraft(nAircraft)

	err := db.BulkInsertAircraftCurrent(aircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err)
	}

	n := 0
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM aircraft_current").Scan(&n)
	if err != nil {
		t.Fatalf("error counting aircraft: %q", err)
	}

	assert.Equal(t, nAircraft, n)
}

func TestAdsbDB_BulkInsertAircraftCurrent_MaxPostgresParameters(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	var maxAircraft = 65535/9 + 1

	aircraft := testUtility.CreateMockAircraft(maxAircraft)

	err := db.BulkInsertAircraftCurrent(aircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err)
	}

	n := 0
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM aircraft_current").Scan(&n)
	if err != nil {
		t.Fatalf("error counting aircraft: %q", err)
	}

	assert.Equal(t, maxAircraft, n)
}

func TestAdsbDB_BulkInsertAircraftCurrent_InvalidType(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	// Create an aircraft with a null icao value
	aircraft := []models.AircraftCurrentModel{
		{
			Icao:         "", // null icao value
			Callsign:     "",
			Altitude:     0,
			Latitude:     0,
			Longitude:    0,
			Speed:        0,
			Track:        0,
			VerticalRate: 0,
			Timestamp:    time.Now().String(),
		},
	}

	err := db.BulkInsertAircraftCurrent(aircraft)

	if err == nil {
		t.Fatalf("Expected an error when inserting invalid data, got nil")
	}
}

func TestAdsbDB_InsertHistoryFromCurrent(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	var nAircraft = 100

	mockAircraft := testUtility.CreateMockAircraft(nAircraft)

	err := db.BulkInsertAircraftCurrent(mockAircraft)
	if err != nil {
		t.Fatalf("error inserting mockAircraft: %q", err)
	}

	err = db.InsertHistoryFromCurrent()
	if err != nil {
		t.Fatalf("error adding history data: %q", err)
	}

	n := 0
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM aircraft_history").Scan(&n)
	if err != nil {
		t.Fatalf("error counting mockAircraft: %q", err)
	}

	assert.Equal(t, nAircraft, n)
}

func TestAdsbDB_DeleteOldCurrent(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	acAfter := testUtility.CreateMockAircraftWithTimestamp("TEST1",
		time.Now().Add(-(global.WaitingTime+3)*time.Second).Format(time.DateTime))

	acNow := testUtility.CreateMockAircraftWithTimestamp("TEST2",
		time.Now().Format(time.DateTime))

	err := db.BulkInsertAircraftCurrent([]models.AircraftCurrentModel{acAfter, acNow})
	if err != nil {
		t.Fatalf("Error inserting aircraft: %q", err)
	}

	err = db.DeleteOldCurrent()
	if err != nil {
		t.Fatalf("Error deleting old aircraft: %q", err)
	}

	var count int

	// check if the old aircraft is deleted
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM aircraft_current WHERE icao = $1", acAfter.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Old aircraft data was not deleted")
	}

	// check if the recent aircraft data is still there
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM aircraft_current WHERE icao = $1", acNow.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 1 {
		t.Fatalf("Recent aircraft data was deleted")
	}
}

func TestAdsbDB_GetCurrentAircraft(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	acAfter := testUtility.CreateMockAircraftWithTimestamp("TEST1",
		time.Now().Add(-(global.WaitingTime+3)*time.Second).Format(time.DateTime))

	var icaoTest2 = "TEST2"
	acNow := testUtility.CreateMockAircraftWithTimestamp(icaoTest2,
		time.Now().Format(time.DateTime))

	var count = 0
	aircraft, err := db.GetCurrentAircraft()
	if err != nil {
		t.Fatalf("Error getting all current aircraft: %q", err)
	}

	count = len(aircraft)
	if count != 0 {
		t.Fatalf("Expected error, db should not contain any elements")
	}

	err = db.BulkInsertAircraftCurrent([]models.AircraftCurrentModel{acAfter, acNow})
	if err != nil {
		t.Fatalf("Error inserting aircraft: %q", err)
	}

	aircraft, err = db.GetCurrentAircraft()
	if err != nil {
		t.Fatalf("Error getting all current aircraft: %q", err)
	}

	count = len(aircraft)

	if count != 1 {
		t.Fatalf("Expected error, list should only contain 1 element")

	}

	assert.Equal(t, icaoTest2, aircraft[0].Icao)
}

func TestAdsbDB_GetHistoryByIcao(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	var nAircraft = 100
	var icao = "TEST"
	mockAircraft := testUtility.CreateMockAircraftWithIcao(nAircraft, icao)

	err := db.BulkInsertAircraftCurrent(mockAircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err.Error())
	}

	err = db.InsertHistoryFromCurrent()
	if err != nil {
		t.Fatalf("error inserting history data: %q", err.Error())
	}

	aircraft, err := db.GetHistoryByIcao(icao)
	if err != nil {
		t.Fatalf("error retriving history data: %q", err.Error())
	}

	assert.Equal(t, nAircraft, len(aircraft))
	for i, ac := range aircraft {
		assert.Equal(t, mockAircraft[i].Icao, ac.Icao)
		assert.Equal(t, mockAircraft[i].Latitude, ac.Latitude)
		assert.Equal(t, mockAircraft[i].Longitude, ac.Longitude)
	}
}

func TestAdsbDB_GetHistoryByIcao_InvalidIcao(t *testing.T) {
	db := setupTestDB()
	defer teardownTestDB(db)

	err := db.InsertHistoryFromCurrent()
	if err != nil {
		t.Fatalf("error inserting history data: %q", err.Error())
	}

	aircraft, err := db.GetHistoryByIcao("")
	if err != nil {
		t.Fatalf("error retriving history data: %q", err.Error())
	}

	assert.Equal(t, 0, len(aircraft))
}
