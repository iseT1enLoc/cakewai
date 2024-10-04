package response

const (
	INTERNAL_SERVER_ERROR = "internal_server_error"
	SUCCESS               = "success"
)

// FailedResponse represents a failed response structure for API responses.
type FailedResponse struct {
	Code    int    `json:"code" example:"500"`                      // HTTP status code.
	Message string `json:"message" example:"internal_server_error"` // Message corresponding to the status code.
	Error   string `json:"error" example:"{$err}"`                  // error message.
}

// BasicResponse represents a failed response structure for API responses.
type BasicResponse struct {
	Code    int         `json:"code" example:"500"`                      // HTTP status code.
	Message string      `json:"message" example:"internal_server_error"` // Message corresponding to the status code.
	Error   string      `json:"error" example:"{$err}"`                  // error message.
	Data    interface{} `json:"data,omitempty"`
}

// BasicBuilder constructs a BasicBuilder based on the provided error.
func BasicBuilder(result BasicResponse) BasicResponse {
	return result
}

// SuccessResponse represents a success response structure for API responses.
type SuccessResponse struct {
	Success
	Meta
}

type ResponseFormat struct {
	Code    int    `json:"code" example:"200"` // HTTP status code.
	Message string `json:"message" example:"success"`
}

type Success struct {
	ResponseFormat
	Data interface{} `json:"data,omitempty"` // data payload.
}

type Meta struct {
	Meta interface{} `json:"meta,omitempty"` //pagination payload.
	Success
}
