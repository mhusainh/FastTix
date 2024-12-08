package response

import (
	"net/http"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Meta: Meta{Code: http.StatusOK, Message: message},
		Data: data,
	}
}

func ErrorResponse(code int, message string) Response {
	return Response{
		Meta: Meta{Code: code, Message: message},
		Data: nil,
	}
}
