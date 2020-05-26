package covidtracker

import "fmt"

// General errors.
const (
	ErrUnauthorized = Error("unauthorized")
)

// HTTP errors
const (
	ErrInvalidTimeJSON            = Error("invalid time format (only RFC3339 format `2006-01-02T15:04:05Z07:00` is accepted)")
	ErrInvalidJSON                = Error("invalid json")
	ErrInternal                   = Error("internal error")
	ErrInvalidQueryParamsType     = Error("invalid query parameters type")
	ErrNoAuthenticationToken      = Error("no authentication token provided")
	ErrInvalidAuthenticationToken = Error("invalid authentication token provided")
	ErrUnauthorizedUser           = Error("user is not authorized for back office actions")
)

//DB Errors
var (
	ErrDocRequired = func(doc string) error {
		return Errorf("document %q is required", doc)
	}
)

// Error represents an error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }

// Create a billing Error
func Errorf(pattern string, params ...interface{}) Error {
	return Error(fmt.Sprintf(pattern, params...))
}
