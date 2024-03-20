package db

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/testUtility"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	global.InitTestEnv()
	m.Run()
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

	err = db.Conn.Close()
	if err != nil {
		t.Errorf("Database connection failed: %q", err)
	}
}

func TestAdsbDB_CreateCurrentTimeAircraftTable(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer teardownTestDB(db)

	err = db.CreateCurrentTimeAircraftTable()
	if err != nil {
		t.Errorf("Create current_time_aircraft table failed: %q", err)
	}

	query := `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'current_time_aircraft')`

	var exists bool
	err = db.Conn.QueryRow(query).Scan(&exists)
	if err != nil {
		t.Fatalf("table does not exists: %q", err)
	}

	query = `SELECT EXISTS (SELECT 1 FROM   pg_indexes WHERE  indexname = $1 AND    tablename = $2)`

	err = db.Conn.QueryRow(query, "timestamp_index", "current_time_aircraft").Scan(&exists)
	if err != nil {
		t.Fatalf("index does not exists: %q", err)
	}
}

func TestAdsbDB_BulkInsertCurrentTimeAircraftTable(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer teardownTestDB(db)

	var nAircraft = 100

	aircraft := testUtility.CreateMockAircraft(nAircraft)

	err = db.BulkInsertCurrentTimeAircraftTable(aircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err)
	}

	n := 0
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM current_time_aircraft").Scan(&n)
	if err != nil {
		t.Fatalf("error counting aircraft: %q", err)
	}

	assert.Equal(t, nAircraft, n)
}

func TestAdsbDB_BulkInsertCurrentTimeAircraftTable_MaxPostgresParameters(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer teardownTestDB(db)

	var maxAircraft = 65535/9 + 1

	aircraft := testUtility.CreateMockAircraft(maxAircraft)

	err = db.BulkInsertCurrentTimeAircraftTable(aircraft)
	if err != nil {
		t.Fatalf("error inserting aircraft: %q", err)
	}

	n := 0
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM current_time_aircraft").Scan(&n)
	if err != nil {
		t.Fatalf("error counting aircraft: %q", err)
	}

	assert.Equal(t, maxAircraft, n)
}

func TestAdsbDB_BulkInsertCurrentTimeAircraftTable_InvalidType(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer teardownTestDB(db)

	// Create an aircraft with a null icao value
	aircraft := []global.Aircraft{
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

	err = db.BulkInsertCurrentTimeAircraftTable(aircraft)

	if err == nil {
		t.Fatalf("Expected an error when inserting invalid data, got nil")
	}
}

func TestAdsbDB_DeleteOldCurrentAircraft(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer teardownTestDB(db)

	acAfter := testUtility.CreateMockAircraftWithTimestamp("TEST1",
		time.Now().Add(-(global.AdsbHubTime+1)*time.Second).Format(time.DateTime))

	acNow := testUtility.CreateMockAircraftWithTimestamp("TEST2",
		time.Now().Format(time.DateTime))

	err = db.BulkInsertCurrentTimeAircraftTable([]global.Aircraft{acAfter, acNow})
	if err != nil {
		t.Fatalf("Error inserting aircraft: %q", err)
	}

	err = db.DeleteOldCurrentAircraft()
	if err != nil {
		t.Fatalf("Error deleting old aircraft: %q", err)
	}

	var count int

	// check if the old aircraft is deleted
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM current_time_aircraft WHERE icao = $1", acAfter.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Old aircraft data was not deleted")
	}

	// check if the recent aircraft data is still there
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM current_time_aircraft WHERE icao = $1", acNow.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 1 {
		t.Fatalf("Recent aircraft data was deleted")
	}
}

func TestAdsbDB_GetAllCurrentAircraft(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer teardownTestDB(db)

	acAfter := testUtility.CreateMockAircraftWithTimestamp("TEST1",
		time.Now().Add(-(global.AdsbHubTime+1)*time.Second).Format(time.DateTime))

	var icaoTest2 = "TEST2"
	acNow := testUtility.CreateMockAircraftWithTimestamp(icaoTest2,
		time.Now().Format(time.DateTime))

	var count = 0
	geoJsonFeatureCollection, err := db.GetAllCurrentAircraft()
	if err != nil {
		t.Fatalf("Error getting all current geoJsonFeatureCollection: %q", err)
	}

	count = len(geoJsonFeatureCollection.Features)
	if count != 0 {
		t.Fatalf("Expected error, db should not contain any elements")
	}

	err = db.BulkInsertCurrentTimeAircraftTable([]global.Aircraft{acAfter, acNow})
	if err != nil {
		t.Fatalf("Error inserting geoJsonFeatureCollection: %q", err)
	}

	geoJsonFeatureCollection, err = db.GetAllCurrentAircraft()
	if err != nil {
		t.Fatalf("Error getting all current geoJsonFeatureCollection: %q", err)
	}

	count = len(geoJsonFeatureCollection.Features)

	if count != 1 {
		t.Fatalf("Expected error, list should only contain 1 element")

	}

	assert.Equal(t, icaoTest2, geoJsonFeatureCollection.Features[0].Properties.Icao)
}

func setupTestDB() (*AdsbDB, error) {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}
	err = db.CreateCurrentTimeAircraftTable()
	if err != nil {
		log.Fatalf("error creating table: %q", err)
	}
	return db, err
}

func teardownTestDB(db *AdsbDB) {
	_, err := db.Conn.Exec("DROP TABLE IF EXISTS current_time_aircraft")
	if err != nil {
		log.Fatalf("error dropping table: %q", err)
	}

	err = db.CreateCurrentTimeAircraftTable()
	if err != nil {
		log.Fatalf("error creating table: %q", err)
	}

	err = db.Close()
	if err != nil {
		log.Fatalf("error closing database: %q", err)
	}
}
