package common

// ErrorResponse standard error response
// swagger:response errorResponse
type ErrorResponse struct {
	// in: body
	Body struct {
		// required: true
		Message string `json:"message,required"`
	}
}
