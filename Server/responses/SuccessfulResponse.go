package responses

type SuccessfulResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

/**
 * Create a new MissingClientIdResponse.
 */
func NewSuccessfulResponse(message string) *SuccessfulResponse {
	return &SuccessfulResponse{
		StatusCode: 200,
		Message:    message,
	}
}
