package error

import "fmt"

// Create base error class for sending requests.
type RequestError struct {
	// HTTP status code.
	StatusCode int

	// Error message.
	Message string

	// Type of HTTP request.
	RequestType string
}

// Create Error string.
func (re *RequestError) Error() string {
	return fmt.Sprintf("Status Code: %d, Request Type: %s, Error Message: %s", re.StatusCode, re.RequestType, re.Message)
}

// Create a Request Error struct.
func NewRequestError(requestType string, errorMessage string, statusCode int) *RequestError {
	return &RequestError{
		StatusCode:  statusCode,
		Message:     errorMessage,
		RequestType: requestType,
	}
}

// Create base error for general error within codebase.
type GeneralError struct {

	// Function error appeared in.
	Function string

	// Error message.
	Message string
}

// Create Error string.
func (ge *GeneralError) Error() string {
	return fmt.Sprintf("Error in Function: %s, Error Message: %s", ge.Function, ge.Message)
}

// Create a General Error struct.
func NewGeneralError(function string, errorMessage string) *GeneralError {
	return &GeneralError{
		Function: function,
		Message:  errorMessage,
	}
}
