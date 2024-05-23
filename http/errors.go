package http

import (
	"fmt"
)

type ApiError struct {
	HttpStatus int
	Code       int
	Message    string
	TraceID    string
}

func (ae *ApiError) Error() string {
	return fmt.Sprintf("longbridge openapi error, httpStatus:%d code:%d message:%s trace:%s", ae.HttpStatus, ae.Code, ae.Message, ae.TraceID)
}

func NewError(httpStatus int, resp *apiResponse) error {
	return &ApiError{
		HttpStatus: httpStatus,
		Code:       resp.Code,
		Message:    resp.Message,
		TraceID:    resp.TraceID,
	}
}
