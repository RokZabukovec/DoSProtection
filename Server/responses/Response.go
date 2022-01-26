package responses

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

/**
 * Create a new MissingClientIdResponse.
 */
func NewResponse(code int, message string) *Response {
	return &Response{
		StatusCode: code,
		Message:    message,
	}
}
