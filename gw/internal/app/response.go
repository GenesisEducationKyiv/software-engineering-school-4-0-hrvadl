package app

// newErrResponse constructs the JSON encoded response,
// which represents request failure from API.
// NOTE: it expects error to be not nil value.
func newErrResponse(msg string) response[*struct{}] {
	return response[*struct{}]{
		Err:     msg,
		Success: false,
	}
}

// newSuccessResponse constructs the JSON encoded response,
// which represents request success from API.
// NOTE: if data is nil then it will be omitter in the result.
func newSuccessResponse[T any](msg string, data T) *response[T] {
	return &response[T]{
		Msg:     msg,
		Success: true,
		Data:    data,
	}
}

// response struct is a JSON encoded response,
// which represents request success from API.
// NOTE: if data is nil then it will be omitter in the result.
type response[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"message,omitempty"`
	Data    T      `json:"data"`
	Err     string `json:"error,omitempty"`
}
