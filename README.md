# ADSB-Reception_Processing_Display_Analysis
ADS-B - Reception, Processing, Display and Analysis. A bachelor thesis by Andreas Follvaag Malde and Fredrik Sundt-Hansen at NTNU Gjøvik.

## Table of Contents
- [Project Description](#project-description)
    - [REST API](#rest-api)
      - [Database](#Database)
    - [SBS Flight Traffic Receiving API](#sbs-receiving-api)
    - [Website in React](#website)
    - [Logging](#logging)
- [Deployment](#deployment)
- [Testing](#testing)
- [License](#license)
- [Contact](#contact)


### Project Description
This GitHub repository contains the application developed for our bachelor's thesis in Computer Science at the 
Norwegian University of Science and Technology. The purpose of this project was to build an application for 
Electronic Chart Centre AS, the product owners, so they could use it in their applications. Or any organization, 
for that matter, could clone/fork this repository to use as they like. 

It consists of a fullstack application able to receive and process live SBS flight traffic from a 
given SBS source, expose that data through a REST API, and show live and historic data on a website. Golang was used for 
the backend, Postgres was used for the database, and React was used for the frontend.

The application is deployed on the university's internal infrastructure:
`http://129.241.150.147/` - website
`http://129.241.150.147:8080/` - RESTful API

#### SBS Receiving API
`backend/cmd/reception/main.go` Consists of an infinite loop that processes SBS data by receiving data through a TCP 
stream from the given SBS source and converts the data to aircraft structs. Then it inserts that  
newly gotten data into the database. Lastly, it runs a cleanup job using crontab to delete old data from database to 
restrict it from getting too big.

When processing the SBS data the program assumes that there is a time-period between each batch of new data, this is 
the WaitingTime variable. 

When processing SBS data there are three outcomes:
1. There is an error connecting to the source. It will then log the error, sleep WaitingTime seconds, and retry. 
2. It successfully connected to the source but received no data. That is, it got no data in the time between a
WaitingTime period. It sleeps and retires. 
3. It successfully connected to the source and received data. It will then continue on adding this data to the database.
At the end it will sleep UpdatingPeriod seconds and do another iteration. 

##### Why an infinite loop?
There is no end condition to the SBS stream we used for developing and testing, `data.adsbhub.org:5002`. 
The SBS source was a continuous stream, and the application was developed with this in mind.
One could change the loop to exit if there is an error connecting the source. However, if there is downtime on their 
side, the whole application would end and one would need to restart it. Thus, we decided to have an infinite loop. 

##### Database
The current database schema does not use any referential integrity constraints, but uses application enforced 
referential integrity. Due to the fact that the relationship between these two tables is 0..1 to 0..*. Since 
aircraft_current only contains the aircraft that are in the air at this current time, new aircraft will not 
have any history yet. On the other hand, aircraft_history might not have any matching aircraft in aircraft_current. 
Due to lack of coverage, aircraft might disappear momentarily and come back again. 

The application enforced referential integrity is handled in `/backend/internal/db/database.go` 

#### REST API
`backend/cmd/rest/main.go` To make the retrieved data available for external resources, such as the website described 
below, a RESTful API has been implemented. 

Endpoints:
````text
/aircraft/current/                                                                                                         
/aircraft/history/
````

##### Current Aircraft
This endpoint retrieves all aircraft in aircraft_current table, that is, all aircraft currently in the air. 

```
Method: GET
Path: /aircraft/current/
Content-Type: application/json 
```

Status code:
```
200: OK
204: No Content. Valid request, but the aircraft with that ICAO does not exists in the database.
400: Bad Request. Not a valid URL, ICAO or hour parameter.
405: Method not allowed. 
414: Request URI too long.
500: Internal Server Error. Returned if the service is unable to respond to the request, and there is something 
wrong with the service.
```

Body: 
Follows GeoJSON standard for a Point: `https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.2`
````text
{
    "type": "FeatureCollection",                                        (string)                                 
    "features": <GeoJSON features>                                      (array)
                [
                    "type": "Feature"                                   (string)
                    "properties": <aircraft_database_model_properties>  (object)
                                    "icao": <aircraft_icao>_code>       (string)
                                    "callsign": <aircraft_callsign>     (string)
                                    "altitude": <aircraft_altitude>     (int)
                                    "speed": <aircraft_speed>           (int)
                                    "track": <aircraft_track>           (int)
                                    "vspeed": <aircraft_vertical_speed> (int)
                                    "timestamp": <aircraft_timestamp>   (string)
                    "geometry": <GeoJSON geometry>                      (object)
                                "type": "Point"                         (string)
                                "coordinates": [                        (array)
                                                    [
                                                      <latitude>,       (float32)
                                                      <longitude>       (float32)
                                                    ]          
                                               ],
                                
                ]   
}
````
Example request: `/aircraft/current/`
Response:
````json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          -26.072605,
          28.25768
        ]
      },
      "properties": {
        "icao": "834D",
        "callsign": "LNK036E",
        "altitude": 6325,
        "speed": 220,
        "track": 16,
        "vspeed": 640,
        "timestamp": "2024-04-11T20:15:08Z"
      }
    },
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          -33.987366,
          25.355324
        ]
      },
      "properties": {
        "icao": "AC43",
        "callsign": "ZUISJ",
        "altitude": 7625,
        "speed": 84,
        "track": 276,
        "vspeed": 640,
        "timestamp": "2024-04-11T20:15:08Z"
      }
    }
  ]
}  
````

##### Aircraft History
This endpoint retrieves the history of one aircraft by searching for its unique ICAO code. 
Additionally, it also has an optional query parameter 'hour' to limit the history result. 

Header: 
```
Method: GET
Path: /aircraft/history/{icao}?hour=
Content-Type: application/json 
```

Status code: 
```
200: OK
204: No Content. Valid request, either there were no history or only instance, point, for that ICAO.
400: Bad Request. Not a valid URL, ICAO or hour parameter.
405: Method not allowed. 
414: Request URI too long.
500: Internal Server Error. Returned if the service is unable to respond to the request, and there is something 
wrong with the service.
```

Body: 
Follows GeoJSON standard for a LineString: `https://datatracker.ietf.org/doc/html/rfc7946#section-3.1.4`

````text
{
    "type": "FeatureCollection",                                        (string)                                 
    "features": <GeoJSON features>                                      (array)
                [
                    "type": "Feature"                                   (string)
                    "properties": <aircraft_database_model_properties>  (object)
                                    "icao": <aircraft_icao>_code>       (string)
                    "geometry": <GeoJSON geometry>                      (object)
                                "coordinates": [                        (array)
                                                    [
                                                      <latitude>,       (float32)
                                                      <longitude>       (float32)
                                                    ]          
                                               ],
                                "type": "LineString"                    (string)
                ]   
}
````
Example request: `/aircraft/history/101BC`                                                                               
Response: 
````json
{
    "type": "FeatureCollection",
    "features": [
        {
            "type": "Feature",
            "properties": {
                "icao": "101BC"
            },
            "geometry": {
                "coordinates": [
                    [
                        1.830139,
                        39.026554
                    ],
                    [
                        1.850181,
                        39.00128
                    ],
                    [
                        1.869885,
                        38.976334
                    ],
                   ...
                ],
                "type": "LineString"
            }
        }
    ]
}
````
Example request: `/aircraft/history/101BC?hour=1`                                                                 
Response: 
````json
{
    "type": "FeatureCollection",
    "features": [
        {
            "type": "Feature",
            "properties": {
                "icao": "101BC"
            },
            "geometry": {
                "coordinates": [
                    [
                        1.830139,
                        39.026554
                    ],
                    [
                        1.850181,
                        39.00128
                    ],
                   ...
                ],
                "type": "LineString"
            }
        }
    ]
}
````

### Logging
For logging, the 'zerolog' library was used `github.com/rs/zerolog`. The global logging level is set by the environment
variable 'ENV.' For production environment: ENV=production, sets the global logging level to Warning, e.i., all logs with 
Warning, Error, Fatal, Panic will be logged. Any other value than 'prod' or 'production' will set the global logging 
level to trace, all levels are logged, Trace, Debug, Info, Warning etc.  

#### Website

## Deployment 
For deploying the service, Docker is used. The Docker Compose file in the root of the project orchestrates deploying of 
backend with frontend. However, if one would like to only deploy the frontend, there is a compose file in the frontend 
directory. As described above, `REACT_APP_SERVER` can be set to a remote URL hosting the backend. 

Additionally, in the project root folder, there is a .env file for setting environment variables. These are variables
 the backend uses for logging into the database, set global values like WaitingTime, CleaningPeriod, etc. 
For developing, default values are set in `backend/internal/global/db.go` for database variables, and 
`backend/internal/global/sbs.go` for SBS variables. An instance of the environment variable in the .env file will 
overwrite the default values. 

Environment variables:
- DB_USER, database username, No default value
- DB_PASSWORD, database password, No default value
- DB_NAME, database name, Default value: adsb_db
- DB_HOST, database host, Default value: localhost
- DB_PORT, database port, Default Value: 5432
- WAITING_TIME, time between each batch of SBS data, Default value: 4 seconds
- CLEANING_PERIOD, crontab schedule for cleaning old data, Default value is once a day: 0 0 * * *
- UPDATING_PERIOD, time between next for-loop iteration, Default value: 10 seconds
- MAX_DAYS_HISTORY, max amount of history to keep in the database, Default value: 1 day
- SBS_SOURCE, URL for the SBS source to be used for retrieving flight data, No default value


## Testing
Throughout the project, a combination of unit testing and integration testing is used. For testing individual database
functions, a system of setup and teardown is used. The function `setupTestDB(t *testing.T) *Context` will create all
necessary tables for testing, regardless if the tables are going to be used or not. Then a defer statement is used for
the function `teardownTestDB(ctx *Context, t *testing.T)`. This will drop all tables and close the connection. 
This is done with all database tests, making them more integration tests than unit tests-. This was implemented instead
of using database mocks.

All other parts of the application utilize unit tests. Tests that require database interactions make use of a mock
database through the 'gomock' library at 'github.com/golang/mock/gomock'.


For GeoJSON testing, the comprehensive GeoJSON schema available at https://geojson.org/schema/GeoJSON.json is utilized.

## License
MIT License

Copyright (c) 2024 Andreas Follevaag Malde

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Contact
fredsu@stud.ntnu.no
andrefma@stud.ntnu.no