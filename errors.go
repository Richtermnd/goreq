package goreq

import "errors"

type GoreqError struct {
	Message  string `json:"message"`
	HttpCode int    `json:"-"`
	Err      error  `json:"-"`
}

func newGoreqError(message string, httpCode int, err error) GoreqError {
	return GoreqError{
		Message:  message,
		HttpCode: httpCode,
		Err:      err,
	}
}

func (e GoreqError) Error() string {
	return e.Message
}

func (e GoreqError) Unwrap() error {
	return e.Err
}

type ErrorResponseLevel int

const (
	QuietLevel    ErrorResponseLevel = 4  // Only {"msg": "Bad request"}
	InfoLevel     ErrorResponseLevel = 0  // With request error details {"msg": "Wrong type foo field"}
	DetailedLevel ErrorResponseLevel = -4 // Raw go error in response
)

// Server errors
var (
	ErrNotPointerToStruct = errors.New("dest must be a pointer to struct")
	ErrInvalidDestination = errors.New("invalid destination")
)

// Client errors
var (
	ErrBadRequest = errors.New("bad request")
)
