package response

type Response struct {
	Status string `json:"status"` // Error, Ok
	Error string `json:"error,omitempty"`
}

const (
	StatusOk = "Ok"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error: msg,
	}
}