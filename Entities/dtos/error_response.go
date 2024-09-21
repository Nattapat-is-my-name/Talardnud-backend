package dtos

type ErrorResponse struct {
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Error message
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
