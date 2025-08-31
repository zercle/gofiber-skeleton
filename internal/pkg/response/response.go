package response

type JSendResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Message string     `json:"message,omitempty"`
}

func Success(data interface{}) *JSendResponse {
	return &JSendResponse{
		Status: "success",
		Data:   data,
	}
}

func Fail(data interface{}) *JSendResponse {
	return &JSendResponse{
		Status: "fail",
		Data:   data,
	}
}

func Error(message string) *JSendResponse {
	return &JSendResponse{
		Status:  "error",
		Message: message,
	}
}