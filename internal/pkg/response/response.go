package response

type JSendResponse struct {
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func Success(data any) *JSendResponse {
	return &JSendResponse{
		Status: "success",
		Data:   data,
	}
}

func Fail(data any) *JSendResponse {
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
