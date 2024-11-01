package werror

import "fmt"

// Compile-time proof
var (
	_ error  = (*Error)(nil)
	_ WError = (*Error)(nil)
)

// Sentinel errors with status codes
var (
	ErrMethodNotAllowed     = New("This code can only handle GET requests", true, 405)
	ErrCityParameterMissing = New("City parameter is missing", false, 400)
	ErrAPIKeyNotFound       = New("API key not found", true, 500)
	ErrRequestFailed        = New("Failed to make request", true, 500)
	ErrParseResponse        = New("Failed to parse response", true, 500)
	ErrLocationDataMissing  = New("Failed to retrieve location data", true, 500)
	ErrWeatherDataMissing   = New("Failed to retrieve current weather data", true, 500)
	ErrConditionDataMissing = New("Failed to retrieve weather condition data", true, 500)
	ErrInvalidCity          = New("Invalid city name. Please try again.", false, 400)
	ErrForbiddenAccess      = New("Forbidden access", true, 403)
	ErrUnexpectedStatus     = New("Unexpected status code", true, 500)
	ErrCharacterLessThan    = New("City value shouldn't be less than 2 characters", false, 400)
	ErrCharacterMoreThan    = New("City value shouldn't be more than 30 characters", false, 400)
)

// WError interface defines custom error behaviors
type WError interface {
	Wrap(err error) WError
	Unwrap() error
	AddData(data any) WError
	ClearData() WError
	Error() string
	Code() int
}

// Error struct holds custom error information
type Error struct {
	Err        error
	Message    string
	Data       any `json:"-"`
	Loggable   bool
	statusCode int
}

// AddData adds extra data to the error instance
func (e *Error) AddData(data any) WError {
	e.Data = data
	return e
}

// Unwrap retrieves the wrapped error
func (e *Error) Unwrap() error {
	return e.Err
}

// ClearData removes any extra data attached to the error
func (e *Error) ClearData() WError {
	e.Data = nil
	return e
}

// Wrap associates a new error with the existing error chain
func (e *Error) Wrap(err error) WError {
	e.Err = err
	return e
}

// Error returns the complete error message, including wrapped errors if present
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s, %s", e.Err.Error(), e.Message)
	}
	return e.Message
}

// Code returns the HTTP status code of the error
func (e *Error) Code() int {
	return e.statusCode
}

// New creates a new WError instance with a message, loggable flag, and status code
func New(message string, loggable bool, statusCode int) WError {
	return &Error{
		Message:    message,
		Loggable:   loggable,
		statusCode: statusCode,
	}
}
