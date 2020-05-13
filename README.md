# API to evaluate risk of COVID-19

## Prerequisite

In order to run the API server, you need a mongodb database

    brew install mongodb
    brew services start mongodb

## Configure

@todo

## Getting started

Build and run the API server with

    cd cmd/covid-tracker
    go build
    ./covid-tracker

or

    go run main.go

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

### Integration tests

@todo
An 'a minima' integration test to create an operation can be executed by running

    go run cmd/integration/main.go

## Deployment

TODO