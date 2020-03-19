package api

import (
	"fmt"
	"net/http"
)

type ErrorTemplate Error

var (
	// System Error
	ErrorUnknown            = ErrorTemplate{Code: 10001, Message: "Unknown Error"}
	ErrorServiceUnavailable = ErrorTemplate{Code: 10002, Message: "Service unavailable"}
	ErrorRemoteService      = ErrorTemplate{Code: 10003, Message: "Remote service error"}
	ErrorRateLimit          = ErrorTemplate{Code: 10004, Message: "Rate limit"}
	ErrorPermissionDenied   = ErrorTemplate{Code: 10005, Message: "Permission denied", status: http.StatusForbidden}
	ErrorIllegalRequest     = ErrorTemplate{Code: 10006, Message: "Illegal request", status: http.StatusBadRequest}
	ErrorHTTPMethod         = ErrorTemplate{Code: 10007, Message: "HTTP method is not suported for this request", status: http.StatusMethodNotAllowed}
	ErrorParameter          = ErrorTemplate{Code: 10008, Message: "Parameter error", status: http.StatusBadRequest}
	ErrorMissParameter      = ErrorTemplate{Code: 10009, Message: "Miss required parameter", status: http.StatusBadRequest}
	ErrorDB                 = ErrorTemplate{Code: 10010, Message: "DB error"}
	ErrorTokenInvaild       = ErrorTemplate{Code: 10011, Message: "Token invaild", status: http.StatusUnauthorized}
	ErrorMissToken          = ErrorTemplate{Code: 10012, Message: "Miss token", status: http.StatusUnauthorized}
	ErrorVersion            = ErrorTemplate{Code: 10013, Message: "API version %s invalid"}
	ErrorNotFound           = ErrorTemplate{Code: 10014, Message: "Resource not found", status: http.StatusNotFound}
	ErrorInvalidFields      = ErrorTemplate{Code: 10015, Message: "Invalid fields [ %v ]", status: http.StatusBadRequest}

	// Product Error
	ErrorInvalidCodeError = ErrorTemplate{Code: 20001, Message: "Code is invalid", status: http.StatusOK}
	ErrorOutOfStockError  = ErrorTemplate{Code: 20002, Message: "Out of stock", status: http.StatusOK}
)

var errorMessagePrefix string

func SetErrorMessagePrefix(s string) {
	errorMessagePrefix = s
}

func (t ErrorTemplate) New(err error, v ...interface{}) *Error {
	e := Error{
		Code:    t.Code,
		Message: fmt.Sprintf(t.Message, v...),
		err:     err,
	}
	if err != nil {
		if errorMessagePrefix == "" {
			e.Details = err.Error()
		} else {
			e.Details = fmt.Sprintf("%s: %s", errorMessagePrefix, err.Error())
		}
	}
	return &e
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Details
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *Error) Status() int {
	if e == nil || e.status == 0 {
		return http.StatusInternalServerError
	}
	return e.status
}
