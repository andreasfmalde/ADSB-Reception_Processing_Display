package global

// Database variables
var (
	User     string
	Password string
	Dbname   = "adsb_db"
)

// Database constants
const (
	Host = "localhost"
	Port = 5432
)

// API constants
const (
	DefaultPort         = "8080"
	VERSION             = "1.0.0"
	DefaultPath         = "/"
	CurrentAircraftPath = "/aircraft/current/"
	WaitingTime         = 4
)
