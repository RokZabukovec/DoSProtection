package responses

type MissingParametersResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Parameter  string `json:"parameter"`
}

/**
 * Create a new MissingClientIdResponse.
 */
func NewMissingParameterResponse(parameter string) *MissingParametersResponse {
	return &MissingParametersResponse{
		StatusCode: 400,
		Message:    "Missing parameter.",
		Parameter:  parameter,
	}
}
