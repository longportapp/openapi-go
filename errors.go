package openapi

import "fmt"


type ApiError struct {
	httpStatus int
	code int
	message string
}

func (ae *ApiError) Error() string {
	return fmt.Sprintf("longbridge openapi error, httpStatus:%d code:%d message:%s", ae.httpStatus, ae.code, ae.message)
}

func NewError() error {

}
