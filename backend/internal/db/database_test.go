package db

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/models"
	"adsb-api/internal/utility/logger"
	"adsb-api/internal/utility/testUtility"
	"fmt"
	reflect "reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	m.Run()
}

func setupTestDB() *Context {
	ctx, err := InitDB()
	if err != nil {
		logger.Error.Fatalf("Failed to initialize service: %v", err)
	}

	err = ctx.CreateAircraftCurrentTable()
	if err != nil {
		logger.Error.Fatalf("error creating current_time_aircraft table: %q", err)
	}

	err = ctx.CreateAircraftHistoryTable()
	if err != nil {
		logger.Error.Fatalf("error creating history_aircraft table: %q", err)
	}

	return ctx
}

func teardownTestDB(ctx *Context) {
	err := ctx.DropAircraftCurrentTable()
	if err != nil {
		logger.Error.Fatalf("error dropping aircraft_current: %q", err)
	}

	_, err = ctx.db.Exec("DROP TABLE IF EXISTS aircraft_history CASCADE")
	if err != nil {
		logger.Error.Fatalf("error droppint current_time_aircraft: %q", err.Error())
	}

	err = ctx.Close()
	if err != nil {
		logger.Error.Fatalf("error closing database: %q", err)
	}
}

func TestInitCloseDB(t *testing.T) {
	ctx, err := InitDB()
	if err != nil {
		t.Fatalf("Database connection failed: %q", err)
	}
	defer func(db *Context) {
		err := db.Close()
		if err != nil {
			t.Fatalf("error closing database: %q", err)
		}
	}(ctx)

	err = ctx.Close()
	if err != nil {
		t.Errorf("Database connection failed: %q", err)
	}
}

func Test_InitDB_InvalidUsername(t *testing.T) {
	global.DbUser = "user"
	ctx, err := InitDB()
	if err == nil {
		t.Error("expected error when changing database username")
	}

	assert.Nil(t, ctx)
	global.DbUser = "test"
}

func TestAdsbDB_CreateAdsbTables(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	err := ctx.DropAircraftCurrentTable()
	if err != nil {
		t.Fatalf("error dropping aircraft_current: %q", err)
	}

	_, err = ctx.db.Exec("DROP TABLE IF EXISTS aircraft_history CASCADE")
	if err != nil {
		logger.Error.Fatalf("error droppint current_time_aircraft: %q", err.Error())
	}

	err = ctx.CreateAircraftCurrentTable()
	if err != nil {
		t.Fatalf("error creating current_time_aircraft table: %q", err)
	}

	err = ctx.CreateAircraftCurrentTimestampIndex()
	if err != nil {
		t.Fatalf("error creating aircraft_current timestamp index: %q", err)
	}

	err = ctx.CreateAircraftHistoryTable()
	if err != nil {
		t.Fatalf("error creating history_aircraft table: %q", err)
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

	checkTableColumns(t, ctx, "aircraft_current", expectedCurrentTimeAircraftColumns)
	checkTableColumns(t, ctx, "aircraft_history", expectedHistoryAircraftColumns)
}

func checkTableColumns(t *testing.T, ctx *Context, tableName string, expectedColumns map[string]string) {
	rows, err := ctx.db.Query(`SELECT column_name, data_type, character_maximum_length FROM information_schema.columns WHERE table_name = $1`, tableName)
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
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	var nAircraft = 100

	aircraft := testUtility.CreateMockAircraft(nAircraft)

	err := ctx.BulkInsertAircraftCurrent(aircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err)
	}

	n := 0
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_current").Scan(&n)
	if err != nil {
		t.Fatalf("error counting aircraft: %q", err)
	}

	assert.Equal(t, nAircraft, n)
}

func TestAdsbDB_BulkInsertAircraftCurrent_MaxPostgresParameters(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	var maxAircraft = 65535/9 + 1

	aircraft := testUtility.CreateMockAircraft(maxAircraft)

	err := ctx.BulkInsertAircraftCurrent(aircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err)
	}

	n := 0
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_current").Scan(&n)
	if err != nil {
		t.Fatalf("error counting aircraft: %q", err)
	}

	assert.Equal(t, maxAircraft, n)
}

func TestAdsbDB_BulkInsertAircraftCurrent_InvalidType(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

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

	err := ctx.BulkInsertAircraftCurrent(aircraft)

	if err == nil {
		t.Fatalf("Expected an error when inserting invalid data, got nil")
	}
}

func TestAdsbDB_InsertHistoryFromCurrent(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	var nAircraft = 100

	mockAircraft := testUtility.CreateMockAircraft(nAircraft)

	err := ctx.BulkInsertAircraftCurrent(mockAircraft)
	if err != nil {
		t.Fatalf("error inserting mockAircraft: %q", err)
	}

	err = ctx.InsertHistoryFromCurrent()
	if err != nil {
		t.Fatalf("error adding history data: %q", err)
	}

	n := 0
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history").Scan(&n)
	if err != nil {
		t.Fatalf("error counting mockAircraft: %q", err)
	}

	assert.Equal(t, nAircraft, n)
}

func TestAdsbDB_SelectAllColumnsAircraftCurrent(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	var nAircraft = 100
	mockAircraft := testUtility.CreateMockAircraft(nAircraft)

	err := ctx.BulkInsertAircraftCurrent(mockAircraft)
	if err != nil {
		t.Fatalf("Error inserting aircraft: %q", err)
	}

	aircraft, err := ctx.SelectAllColumnsAircraftCurrent()
	if err != nil {
		t.Fatalf("Error getting all current aircraft: %q", err)
	}

	assert.Equal(t, nAircraft, len(aircraft))

}

func TestAdsbDB_SelectAllColumnHistoryByIcao(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	var nAircraft = 100
	var icao = "TEST"
	mockAircraft := testUtility.CreateMockAircraftWithIcao(nAircraft, icao)

	err := ctx.BulkInsertAircraftCurrent(mockAircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err.Error())
	}

	err = ctx.InsertHistoryFromCurrent()
	if err != nil {
		t.Fatalf("error inserting history data: %q", err.Error())
	}

	aircraft, err := ctx.SelectAllColumnHistoryByIcao(icao)
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

func TestAdsbDB_SelectAllColumnHistoryByIcao_InvalidIcao(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	err := ctx.InsertHistoryFromCurrent()
	if err != nil {
		t.Fatalf("error inserting history data: %q", err.Error())
	}

	aircraft, err := ctx.SelectAllColumnHistoryByIcao("")
	if err != nil {
		t.Fatalf("error retriving history data: %q", err.Error())
	}

	assert.Equal(t, 0, len(aircraft))
}

func TestAdsbDB_DeleteOldHistory(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	ac1 := testUtility.CreateMockAircraftWithTimestamp("TEST1",
		time.Now().Add(-(global.MaxDaysHistory+1)*24*time.Hour).Format(time.DateTime))

	ac2 := testUtility.CreateMockAircraftWithTimestamp("TEST2",
		time.Now().Add(-(global.MaxDaysHistory)*24*time.Hour).Format(time.DateTime))

	ac3 := testUtility.CreateMockAircraftWithTimestamp("TEST3",
		time.Now().Add(-(global.MaxDaysHistory-1)*24*time.Hour).Format(time.DateTime))

	_, err := ctx.db.Exec(`
		INSERT INTO aircraft_history 
		VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12)`,
		ac1.Icao, ac1.Latitude, ac1.Longitude, ac1.Timestamp,
		ac2.Icao, ac2.Latitude, ac2.Longitude, ac2.Timestamp,
		ac3.Icao, ac3.Latitude, ac3.Longitude, ac3.Timestamp)

	var count int
	// check if the old aircraft is deleted
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history").Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 3 {
		t.Fatalf("aircraft was not inserted correctly")
	}

	err = ctx.DeleteOldHistory(global.MaxDaysHistory)
	if err != nil {
		t.Fatalf("Error deleting old aircraft: %q", err)
	}

	// check if the old aircraft is deleted
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history WHERE icao = $1", ac1.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Old aircraft data was not deleted")
	}

	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history WHERE icao = $1", ac2.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Old aircraft data was not deleted")
	}

	// check if the recent aircraft data is still there
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history WHERE icao = $1", ac3.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 1 {
		t.Fatalf("Recent aircraft data was deleted")
	}
}

func TestAdsbDB_TestBegin(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	assert.Nil(t, ctx.tx)

	err := ctx.Begin()
	if err != nil {
		t.Fatalf("Error starting transaction: %q", err)
	}

	assert.NotNil(t, ctx.tx)

	// sets tx to nil so it does not affect any other test methods
	ctx.tx = nil
	assert.Nil(t, ctx.tx)
}

func TestAdsbDB_TestCommit(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	ac := testUtility.CreateMockAircraftWithTimestamp("TEST",
		time.Now().Format(time.DateTime))

	err := ctx.Begin()
	if err != nil {
		t.Fatalf("Error starting transaction: %q", err)
	}

	assert.NotNil(t, ctx.tx)

	// transaction execution
	res, err := ctx.tx.Exec("INSERT INTO aircraft_current VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		ac.Icao, ac.Callsign, ac.Altitude, ac.Latitude, ac.Longitude, ac.Speed, ac.Track, ac.VerticalRate, ac.Timestamp)
	if err != nil {
		t.Fatalf("Error inserting data: %q", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("error getting rows affected: %q", err)
	}

	if affected != 1 {
		t.Fatalf("aircraft was not inserted correctly")
	}

	err = ctx.Commit()
	if err != nil {
		t.Fatalf("Error committing transaction: %q", err)
	}

	assert.Nil(t, ctx.tx)

	var count int
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_current WHERE icao = $1", ac.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 1 {
		t.Fatalf("Data was not inserted")
	}
}

func TestAdsbDB_TestRollback(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	ac := testUtility.CreateMockAircraftWithTimestamp("TEST",
		time.Now().Format(time.DateTime))

	err := ctx.Begin()
	if err != nil {
		t.Fatalf("Error starting transaction: %q", err)
	}

	assert.NotNil(t, ctx.tx)

	res, err := ctx.tx.Exec("INSERT INTO aircraft_current VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		ac.Icao, ac.Callsign, ac.Altitude, ac.Latitude, ac.Longitude, ac.Speed, ac.Track, ac.VerticalRate, ac.Timestamp)
	if err != nil {
		t.Fatalf("Error inserting data: %q", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("error getting rows affected: %q", err)
	}

	if affected != 1 {
		t.Fatalf("aircraft was not inserted correctly")
	}

	err = ctx.Rollback()
	if err != nil {
		t.Fatalf("Error rolling back transaction: %q", err)
	}

	assert.Nil(t, ctx.tx)

	var count int
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_current WHERE icao = $1", ac.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Data was inserted")
	}
}

func TestContext_SelectAllColumnHistoryByIcaoFilterByTimestamp(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	var nAircraft = 10
	var icao = "TEST"

	var mockAircraft []models.AircraftCurrentModel

	// creates TEST aircraft with an hour between each instance
	for i := 0; i < nAircraft; i++ {
		ac := testUtility.CreateMockAircraftWithTimestamp("TEST", time.Now().Add(-time.Duration(i)*time.Hour).Format(time.DateTime))
		mockAircraft = append(mockAircraft, ac)
	}

	for _, ac := range mockAircraft {
		_, err := ctx.db.Exec("INSERT INTO aircraft_history VALUES ($1, $2, $3, $4)",
			ac.Icao, ac.Latitude, ac.Longitude, ac.Timestamp)
		if err != nil {
			t.Fatalf("Error inserting data: %q", err)
		}
	}

	// Selects half of the rows
	aircraft, err := ctx.SelectAllColumnHistoryByIcaoFilterByTimestamp(icao, nAircraft/2)
	if err != nil {
		t.Fatalf("error retriving history data: %q", err.Error())
	}

	assert.Equal(t, nAircraft/2, len(aircraft))
	for i, ac := range aircraft {
		assert.Equal(t, mockAircraft[i].Icao, ac.Icao)
		assert.Equal(t, mockAircraft[i].Latitude, ac.Latitude)
		assert.Equal(t, mockAircraft[i].Longitude, ac.Longitude)
	}
}

func TestContext_SelectAllColumnHistoryByIcaoFilterByTimestamp_NoHistory(t *testing.T) {
	ctx := setupTestDB()
	defer teardownTestDB(ctx)

	var nAircraft = 10
	var icao = "TEST"

	var mockAircraft []models.AircraftCurrentModel

	// creates TEST aircraft with an hour between each instance
	for i := 0; i < nAircraft; i++ {
		ac := testUtility.CreateMockAircraftWithTimestamp("TEST", time.Now().Add(-time.Duration(i)*time.Hour).Format(time.DateTime))
		mockAircraft = append(mockAircraft, ac)
	}

	for _, ac := range mockAircraft {
		_, err := ctx.db.Exec("INSERT INTO aircraft_history VALUES ($1, $2, $3, $4)",
			ac.Icao, ac.Latitude, ac.Longitude, ac.Timestamp)
		if err != nil {
			t.Fatalf("Error inserting data: %q", err)
		}
	}

	// Selects half of the rows
	aircraft, err := ctx.SelectAllColumnHistoryByIcaoFilterByTimestamp(icao, 0)
	if err != nil {
		t.Fatalf("error retriving history data: %q", err.Error())
	}

	assert.Equal(t, 0, len(aircraft))
}
