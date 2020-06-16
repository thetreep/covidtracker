`covidtracker`, also known as EviteCovid (https://evitecovid.fr/) is an API used to evaluate risk of COVID-19 contamination during a multi-segment trip (trip combining several mode of transportation or hotel).

It has been built with the objective to evaluate risk of trip within France as a first step. As a result, it relies on the opendata provided by French public health agency 'Santé publique France' (In particular this dataset https://www.data.gouv.fr/fr/datasets/donnees-hospitalieres-relatives-a-lepidemie-de-covid-19/). For the risk about hotels, it relies on the data provided by the hotel aggregator CDS (https://www.cdsgroupe.com/). Also, risk parameters have been completed while having in mind the measures currently applied by carriers, in France.

NB : this project is highly experimental and computed risk is not intended to be accurate (as it depends of a lot of unknown parameters). However it aims to give a basic idea of how risky can be a trip compared with another one.

## Prerequisite

In order to run the API server, you need a mongodb database

    brew install mongodb
    brew services start mongodb

## Configure

To run the API, you can configure the following env variables :

    - THETREEP_COVIDTRACKER_MONGO_URL : the URI to connect to mongo database
    - THETREEP_COVIDTRACKER_DATABASE : the name of the database to connect to 
    - THETREEP_COVIDTRACKER_MONGO_USER : optionally, the user to connect to mongo database
    - THETREEP_COVIDTRACKER_MONGO_PASSWORD : optionally, the password to connect to mongo database
	- THETREEP_COVIDTRACKER_SECRET : the secret used to authenticate frontend requests (static for now)
	- THETREEP_COVIDTRACKER_CDS_API_USER : the user to connect to the API of CDS (hotel aggregator that has built an inventory of the measures against Covid-19 in hotels)
	- THETREEP_COVIDTRACKER_CDS_API_PASSWORD : the password to connect to the API of CDS
	- THETREEP_COVIDTRACKER_CDS_API_DUTY_CODE : the duty code to connect to the API of CDS
    - THETREEP_COVIDTRACKER_LOG_LEVEL: the logging level, can be one of "fatal", "error", "warn", "info", "debug" (logs output is in JSON)

## Getting started

Build and run the API server with

    cd cmd/covidtracker
    go build
    ./covidtracker

or

    go run main.go

You also may need to run the refresher (cron task that refreshes once a day the French emergency open data about Covid-19) with

    go run cmd/refresher/main.go

## Dependencies

Dependencies are automatically managed with `go mod` (requires go >= 1.11)
You just need to export the appropriate variable

    export GO111MODULE=on

Then dependencies will automatically be added when using the `go` toolchain (`go build`, `go run`, `go test`...)

## Testing

### Unit tests

Go unit tests can be run with

    go test ./...

or with all flags (to display all output, test coverage and detect race conditions):

    go test -cover -race -v ./...

