package global

const (
	MethodNotSupported              = "method %s is not supported"
	ErrorRetrievingCurrentAircraft  = "error retrieving current aircraft from database"
	ErrorTongURL                    = "requested URL is too long"
	ErrorInvalidQueryParams         = "invalid query parameters: Endpoint only supports the given parameters: "
	ErrorRetrievingAircraftWithIcao = "error retrieving aircraft history with icao: "
	ErrorConvertingDataToGeoJson    = "error converting aircraft data to Geo Json"
	ErrorGeoJsonTooFewCoordinates   = "coordinates array must have at least 2 items"
	ErrorEncodingJsonData           = "error encoding json data"
)
