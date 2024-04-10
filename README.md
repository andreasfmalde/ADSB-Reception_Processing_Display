# ADSB-Reception_Processing_Display_Analysis
ADS-B - Reception, Processing, Display and Analysis. A bachelor thesis for Andreas Follvaag Malde and Fredrik Sundt-Hansen at NTNU Gjøvik.

## Table of Contents
- [Project Description](#project-description)
    - [REST API](#rest-api)
      - [Database](#Database-)
    - [SBS Flight Traffic Receiving API](#sbs-receiving-api)
    - [Website in React](#website)
- [Development](#development)
    - [Whole Application](#whole-application)
    - [Only Frontend](#only-frontend)
- [Testing](#testing)
    - [Database Testing](#rest-api-testing)
    - [GeoJson-Testing](#sbs-api-testing)
    - [React Website User Testing](#react-website-testing)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)


### Project Description
This is the GitHub repository for the bachelor's thesis at the Norwegian University of Science and Technology 
on Computer science. The purpose of this project was to build an application for Electronic Chart Centre AS, the product owner, so they could 
use it in their applications. Or any organization, for that matter, could clone/fork this repository to use as they like. 

It consists of a fullstack application able to receive and process live SBS flight traffic from a 
given SBS source, expose that data through a REST API, and show live and historic data on a website. Golang was used for 
the backend, Postgres was used for the database, and React was used for the frontend.

#### SBS Receiving API
`backend/cmd/reception/main.go` Consists of an infinite loop that processes SBS data by receiving data through a TCP 
stream from the given SBS source and converts the data to aircraft structs. Then it inserts that  
newly gotten data into the database. Finally, it deletes old data from database to restrict it from getting too big.

Relevant environment variables:
- WaitingTime, time between each batch of SBS data. Default value: 4 
- CleaningPeriod, time between each period cleaning cycle. Default value: 120
- UpdatingPeriod, time between next for-loop iteration. Default value: 10
- MaxDaysHistory, max amount of history to keep in the database. Default value: 1

When processing the SBS data the program assumes that there is a time-period between each batch of new data, this is 
the WaitingTime variable. 

When processing SBS data the are three outcomes
1. There is an error connecting to the source. It will then log the error, sleep WaitingTime seconds, and retry. 
2. It successfully connected to the source but recieved no data. That is, it got no data in the time between a
WaitingTime period. 
3. It succesfully connected to the source and recieved data. It will to continue to add this data to the database.

##### Why an infinite loop?
There is no end condition to the SBS stream we used for developing and testing, `data.adsbhub.org:5002`. 
The application is meant for and developed with this in mind. One could change the loop to exit if there is an error 
connecting the source. However, if there is downtime on their side, the whole application would end and one would need 
to restart it. Thus, we decided to have an inifinte loop. 

##### Database 



#### REST API



#### Endpoints
`
/aircraft/current/
/aircraft/history/
`

#### Current Endpoint

#### History Endpoint

### Logger


#### Website

