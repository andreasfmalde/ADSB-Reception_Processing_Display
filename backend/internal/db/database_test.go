package db

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/global/models"
	"adsb-api/internal/utility/testUtility"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	global.InitTestEnvironment()
	m.Run()
}

func setupTestDB(t *testing.T) *Context {
	ctx, err := InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}

	err = ctx.CreateAircraftCurrentTable()
	if err != nil {
		t.Fatalf("error creating current_time_aircraft table: %q", err)
	}

	err = ctx.CreateAircraftHistoryTable()
	if err != nil {
		t.Fatalf("error creating history_aircraft table: %q", err)
	}

	err = ctx.CreateAircraftHistoryTimestampIndex()
	if err != nil {
		t.Fatalf("error creating timestamp_index: %q", err)
	}

	return ctx
}

func teardownTestDB(ctx *Context, t *testing.T) {
	err := ctx.DropAircraftCurrentTable()
	if err != nil {
		t.Fatalf("error dropping aircraft_current: %q", err)
	}

	_, err = ctx.db.Exec("DROP TABLE IF EXISTS aircraft_history CASCADE")
	if err != nil {
		t.Fatalf("error dropping current_time_aircraft: %q", err.Error())
	}

	if ctx.tx != nil {
		err = ctx.Commit()
		if err != nil {
			t.Fatalf("error committing uncommitted transaction: %q", err)
		}
	}

	err = ctx.Close()
	if err != nil {
		t.Fatalf("error closing database: %q", err)
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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

	// Drop tables that are to be tested
	err := ctx.DropAircraftCurrentTable()
	if err != nil {
		t.Fatalf("error dropping aircraft_current: %q", err)
	}

	_, err = ctx.db.Exec("DROP TABLE aircraft_history CASCADE")
	if err != nil {
		t.Fatalf("error droppint current_time_aircraft: %q", err.Error())
	}

	// Test create tables and indexes
	err = ctx.CreateAircraftCurrentTable()
	if err != nil {
		t.Fatalf("error creating current_time_aircraft table: %q", err)
	}

	err = ctx.CreateAircraftHistoryTable()
	if err != nil {
		t.Fatalf("error creating history_aircraft table: %q", err)
	}

	err = ctx.CreateAircraftHistoryTimestampIndex()
	if err != nil {
		t.Fatalf("error creating aircraft_current timestamp index: %q", err)
	}

	var exists bool

	query := `SELECT EXISTS (
		SELECT 1
		FROM   pg_indexes 
		WHERE  tablename = $1
		AND    indexname = $2
	);`

	err = ctx.db.QueryRow(query, "aircraft_history", "timestamp_index").Scan(&exists)
	if err != nil {
		t.Fatalf("timestmap_index was not created")
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
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			t.Fatalf("error closing rows")
		}
	}(rows)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

	icao := "TEST"
	var nAircraft = 100
	mockAircraft := testUtility.CreateMockHistAircraftWithIcao(nAircraft, icao)

	for _, ac := range mockAircraft {
		_, err := ctx.db.Exec("INSERT INTO aircraft_history VALUES ($1, $2, $3, $4)",
			ac.Icao, ac.Latitude, ac.Longitude, ac.Timestamp)
		if err != nil {
			t.Fatalf("Error inserting data: %q", err)
		}
	}

	aircraft, err := ctx.SelectAllColumnHistoryByIcao(icao)
	if err != nil {
		t.Fatalf("error retriving history data: %q", err.Error())
	}

	assert.Equal(t, nAircraft, len(aircraft))
	for i, ac := range aircraft {
		assert.Equal(t, mockAircraft[i].Icao, ac.Icao)
	}
}

func TestAdsbDB_SelectAllColumnHistoryByIcao_InvalidIcao(t *testing.T) {
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

	// recent aircraft
	ac1 := testUtility.CreateMockAircraftWithTimestamp("TEST1",
		time.Now().Add(-(time.Duration(global.MaxDaysHistory)-1)*24*time.Hour).Truncate(time.Hour).Format(time.DateTime))

	// old aircraft
	ac2 := testUtility.CreateMockAircraftWithTimestamp("TEST2",
		time.Now().Add(-time.Duration(global.MaxDaysHistory+1)*24*time.Hour).Truncate(time.Hour).Format(time.DateTime))

	// old aircraft
	ac3 := testUtility.CreateMockAircraftWithTimestamp("TEST3",
		time.Now().Add(-(time.Duration(global.MaxDaysHistory)+2)*24*time.Hour).Truncate(time.Hour).Format(time.DateTime))

	_, err := ctx.db.Exec(`
		INSERT INTO aircraft_history 
		VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12)`,
		ac1.Icao, ac1.Latitude, ac1.Longitude, ac1.Timestamp,
		ac2.Icao, ac2.Latitude, ac2.Longitude, ac2.Timestamp,
		ac3.Icao, ac3.Latitude, ac3.Longitude, ac3.Timestamp)
	if err != nil {
		t.Fatalf("error inserting test data: %v", err)
	}

	var count int

	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history").Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 3 {
		t.Fatalf("aircraft was not inserted correctly")
	}

	err = ctx.DeleteOldHistory(1)
	if err != nil {
		t.Fatalf("Error deleting old aircraft: %q", err)
	}

	// check if the recent aircraft still exists
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history WHERE icao = $1", ac1.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 1 {
		t.Fatalf("Recent aircraft was deleted")
	}

	// check if old aircraft was deleted
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history WHERE icao = $1", ac2.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Old aircraft data was not deleted")
	}

	// check if old aircraft was deleted
	err = ctx.db.QueryRow("SELECT COUNT(*) FROM aircraft_history WHERE icao = $1", ac3.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Old aircraft data was not deleted")
	}
}

func TestAdsbDB_TestBegin(t *testing.T) {
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

	assert.Nil(t, ctx.tx)

	err := ctx.Begin()
	if err != nil {
		t.Fatalf("Error starting transaction: %q", err)
	}

	assert.NotNil(t, ctx.tx)
}

func TestAdsbDB_TestCommit(t *testing.T) {
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

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

func TestContext_Begin_TxAlreadyInitialized(t *testing.T) {
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

	err := ctx.Begin()
	if err != nil {
		t.Fatalf("error beginning transaction: %q", err)
	}

	err = ctx.Begin()
	assert.Error(t, err, "expected error when beginning transaction and not committed or rolled back")
	assert.Equal(t, errorMsg.TransactionInProgress, err.Error())
}

func TestContext_Commit_TxNotInitialized(t *testing.T) {
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

	err := ctx.Commit()
	assert.Error(t, err, "expected error when committing without initialized transaction")
	assert.Equal(t, errorMsg.NoTransactionInProgress, err.Error())
}

func TestContext_Rollback_TxAlreadyCommitted(t *testing.T) {
	ctx := setupTestDB(t)
	defer teardownTestDB(ctx, t)

	err := ctx.Rollback()
	assert.Error(t, err, "expected error when using rollback without initialized transaction")
	assert.Equal(t, errorMsg.NoTransactionInProgress, err.Error())
}
