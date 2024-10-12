package errors

import (
	"fmt"
)

type BusinessError struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

func NewBusinessError(code int, msg string, data ...interface{}) *BusinessError {
	d := interface{}(nil)
	if len(data) > 0 {
		d = data[0]
	}
	return &BusinessError{Code: code, Msg: msg, Data: d}
}
