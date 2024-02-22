package global

/*
Struct used to represent an aircraft record
*/
type Aircraft struct {
	Icao         string  `json:"icao"`
	Callsign     string  `json:"callsign"`
	Altitude     int     `json:"altitude"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	Speed        int     `json:"speed"`
	Track        int     `json:"track"`
	VerticalRate int     `json:"vspeed"`
	Timestamp    string  `json:"timestamp"`
}
