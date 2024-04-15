package errorMsg

const (
	MethodNotSupported              = "method %s is not supported"
	ErrorRetrievingCurrentAircraft  = "error retrieving current aircraft from database"
	ErrorTongURL                    = "requested URL is too long"
	ErrorInvalidQueryParams         = "invalid query parameter: Endpoint only supports the given parameters"
	ErrorRetrievingAircraftWithIcao = "error retrieving aircraft history with icao"
	ErrorConvertingDataToGeoJson    = "error converting aircraft data to Geo Json"
	ErrorGeoJsonTooFewCoordinates   = "coordinates array must have at least 2 items"
	ErrorEncodingJsonData           = "error encoding json data"
	ErrorClosingDatabase            = "error closing database"
	ErrorCreatingDatabaseTables     = "error creating database tables"
	ErrorInsertingNewSbsData        = "could not insert new SBS data"
	ErrorCouldNotConnectToTcpStream = "could not connect to TCP stream"
	EmptyIcao                       = "ICAO code cannot be empty"
	InvalidQueryParameterHour       = "query parameter 'hour', can only be an integer"
	TransactionInProgress           = "transaction already in progress"
	NoTransactionInProgress         = "no transaction in progress"
	TooLongIcao                     = "ICAO code cannot be longer than 6 characters"
	ErrorDeletingOldHistory         = "error deleting old history"

	InfoOldHistoryDataDeleted = "old history data deleted"
)
