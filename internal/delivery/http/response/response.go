package response

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func Error(message string, err error) Response {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}

	return Response{
		Status:  "error",
		Message: message,
		Error:   errorMessage,
	}
}
