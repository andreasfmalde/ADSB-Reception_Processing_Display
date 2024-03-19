package db

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/test"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	global.InitDevEnv()
	m.Run()
}

func InitTestDB() (*AdsbDB, error) {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	_, err = db.Conn.Exec("DROP TABLE IF EXISTS current_time_aircraft")
	if err != nil {
		return nil, fmt.Errorf("error dropping current_time_aircraft table: %w", err)
	}

	err = db.CreateCurrentTimeAircraftTable()
	if err != nil {
		return nil, fmt.Errorf("error creating current_time_aircraft table: %w", err)
	}

	return db, nil
}

func TestAdsbDB(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (*AdsbDB, error)
		cleanup func(db *AdsbDB)
		wantErr bool
	}{
		{
			name: "valid case",
			setup: func() (*AdsbDB, error) {
				db, err := InitTestDB()
				if err != nil {
					t.Fatalf("Failed to initialize test database: %v", err)
				}
				return db, nil
			},
			cleanup: func(db *AdsbDB) {
				db.Close()
			},
			wantErr: false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.setup()
			if err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}
			defer tt.cleanup(db)

			err = db.CreateCurrentTimeAircraftTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCurrentTimeAircraftTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitCloseDB(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Errorf("Database connection failed: %q", err)
	}
	defer db.Close()

	err = db.Conn.Close()
	if err != nil {
		t.Errorf("Database connection failed: %q", err)
	}

}

func TestAdsbDB_CreateCurrentTimeAircraftTable(t *testing.T) {
	db, err := InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer db.Close()

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

func TestAdsbDB_ValidBulkInsertCurrentTimeAircraftTable(t *testing.T) {
	db, err := InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer db.Close()

	_, err = db.Conn.Exec("DELETE FROM current_time_aircraft")
	if err != nil {
		t.Fatalf("could not reset database")
	}

	var nAircraft = 100

	aircraft := test.GetAircraft(nAircraft)

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

func TestAdsbDB_InvalidBulkInsertCurrentTimeAircraftTable(t *testing.T) {
	db, err := InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer db.Close()

	// Create an aircraft with a null icao value
	aircraft := []global.Aircraft{
		{
			Icao:         "", // null icao value
			Callsign:     "TEST",
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
	db, err := InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer db.Close()

	aircraft1 := global.Aircraft{
		Icao:         "TEST1",
		Callsign:     "TEST",
		Altitude:     10000,
		Latitude:     51.5074,
		Longitude:    0.1278,
		Speed:        450,
		Track:        180,
		VerticalRate: 0,
		Timestamp:    time.Now().Add(-7 * time.Second).Format(time.RFC3339Nano), // 7 seconds ago
	}
	aircraft2 := global.Aircraft{
		Icao:         "TEST2",
		Callsign:     "TEST",
		Altitude:     10000,
		Latitude:     51.5074,
		Longitude:    0.1278,
		Speed:        450,
		Track:        180,
		VerticalRate: 0,
		Timestamp:    time.Now().Format(time.RFC3339Nano), // current time
	}
	err = db.BulkInsertCurrentTimeAircraftTable([]global.Aircraft{aircraft1, aircraft2})
	if err != nil {
		t.Fatalf("Error inserting aircraft: %q", err)
	}

	// Call the DeleteOldCurrentAircraft method
	err = db.DeleteOldCurrentAircraft()
	if err != nil {
		t.Fatalf("Error deleting old aircraft: %q", err)
	}

	// Query the table to check if the old aircraft data has been deleted
	var count int
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM current_time_aircraft WHERE icao = $1", aircraft1.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 0 {
		t.Fatalf("Old aircraft data was not deleted")
	}

	// Check if the recent aircraft data is still there
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM current_time_aircraft WHERE icao = $1", aircraft2.Icao).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying the table: %q", err)
	}
	if count != 1 {
		t.Fatalf("Recent aircraft data was deleted")
	}
}

func TestAdsbDB_GetAllCurrentAircraft(t *testing.T) {
	db, err := InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer db.Close()

	// Insert some aircraft data with different timestamps
	aircraft1 := global.Aircraft{
		Icao:         "TEST1",
		Callsign:     "TEST",
		Altitude:     10000,
		Latitude:     51.5074,
		Longitude:    0.1278,
		Speed:        450,
		Track:        180,
		VerticalRate: 0,
		Timestamp:    time.Now().Add(-7 * time.Second).Format(time.RFC3339Nano), // 7 seconds ago
	}
	aircraft2 := global.Aircraft{
		Icao:         "TEST2",
		Callsign:     "TEST",
		Altitude:     10000,
		Latitude:     51.5074,
		Longitude:    0.1278,
		Speed:        450,
		Track:        180,
		VerticalRate: 0,
		Timestamp:    time.Now().Format(time.RFC3339Nano), // current time
	}
	err = db.BulkInsertCurrentTimeAircraftTable([]global.Aircraft{aircraft1, aircraft2})
	if err != nil {
		t.Fatalf("Error inserting aircraft: %q", err)
	}

	// Call the GetAllCurrentAircraft method
	aircraft, err := db.GetAllCurrentAircraft()
	if err != nil {
		t.Fatalf("Error getting all current aircraft: %q", err)
	}

	// Check if the returned aircraft data is within the current timeframe
	for _, feature := range aircraft.Features {
		timestamp, err := time.Parse(time.RFC3339Nano, feature.Properties.Timestamp)
		if err != nil {
			t.Fatalf("Error parsing timestamp: %q", err)
		}
		if timestamp.Before(time.Now().Add(-6 * time.Second)) {
			t.Fatalf("Returned aircraft data is not within the current timeframe")
		}
	}
}

func TestAdsbDB_GetAllCurrentAircraft_Failure(t *testing.T) {
	db, err := InitTestDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	defer db.Close()

	// Close the database connection to force db.Conn.Query to fail
	db.Close()

	// Call the GetAllCurrentAircraft method
	_, err = db.GetAllCurrentAircraft()
	if err == nil {
		t.Fatalf("Expected an error when calling GetAllCurrentAircraft with a closed database connection, got nil")
	}

	// Reopen the database connection
	db, err = InitDB()
	if err != nil {
		t.Fatalf("Database connection failed: %q", err)
	}

	// Insert some aircraft data with an invalid timestamp to force rows.Scan to fail
	aircraft := global.Aircraft{
		Icao:         "TEST",
		Callsign:     "TEST",
		Altitude:     10000,
		Latitude:     51.5074,
		Longitude:    0.1278,
		Speed:        450,
		Track:        180,
		VerticalRate: 0,
		Timestamp:    "invalid timestamp", // invalid timestamp
	}
	err = db.BulkInsertCurrentTimeAircraftTable([]global.Aircraft{aircraft})
	if err == nil {
		t.Fatalf("Error inserting aircraft: %q", err)
	}

	// Call the GetAllCurrentAircraft method
	_, err = db.GetAllCurrentAircraft()
	if err != nil {
		t.Fatalf("Expected an error when calling GetAllCurrentAircraft with invalid data, got nil")
	}
}