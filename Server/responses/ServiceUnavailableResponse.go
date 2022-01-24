package responses

type ServiceUnavailableResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

/**
 * Create a new MissingClientIdResponse.
 */
func NewServiceUnavailableResponse(message string) *ServiceUnavailableResponse {
	return &ServiceUnavailableResponse{
		StatusCode: 503,
		Message:    message,
	}
}
