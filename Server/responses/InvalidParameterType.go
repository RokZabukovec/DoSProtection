package responses

type InvalidParameterType struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Parameter  string `json:"parameter"`
}

/**
 * Create a new MissingClientIdResponse.
 */
func NewInvalidParameterTypeResponse(parameter string, message string) *InvalidParameterType {
	return &InvalidParameterType{
		StatusCode: 400,
		Message:    message,
		Parameter:  parameter,
	}
}
