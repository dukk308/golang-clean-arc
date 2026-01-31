package common

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func NewResponseSuccess(data any) *Response {
	return &Response{
		Data:    data,
		Message: "OK",
	}
}
