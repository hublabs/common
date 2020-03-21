package api

import (
	"errors"
	"testing"

	"github.com/pangpanglabs/goutils/test"
)

func TestErrorsTest(t *testing.T) {
	t.Run("Internal", func(t *testing.T) {
		err := errors.New("invalid sql")
		err = ErrorUnknown.New(err)

		test.Equals(t, err.Error(), "unknown service: invalid sql")

		var apiError *Error
		ok := errors.As(err, &apiError)
		test.Equals(t, ok, true)
		test.Equals(t, apiError.Code, ErrorUnknown.Code)
		test.Equals(t, apiError.Message, ErrorUnknown.Message)
		test.Equals(t, apiError.Details, "unknown service: invalid sql")

		// Test do not cover
		err = ErrorDB.New(err)
		ok = errors.As(err, &apiError)
		test.Equals(t, ok, true)
		test.Equals(t, apiError.Code, ErrorUnknown.Code)
		test.Equals(t, apiError.Message, ErrorUnknown.Message)
		test.Equals(t, apiError.Details, "unknown service: invalid sql")
	})

	t.Run("External", func(t *testing.T) {
		SetErrorMessagePrefix("serviceB")
		v := &Error{
			Code:    ErrorDB.Code,
			Message: ErrorDB.Message,
			Details: "serviceA: invalid sql",
		}

		err := error(v)
		err = ErrorRemoteService.New(err)

		test.Equals(t, err.Error(), "serviceB: serviceA: invalid sql")

		var apiError *Error
		ok := errors.As(err, &apiError)
		test.Equals(t, ok, true)
		test.Equals(t, apiError.Code, ErrorRemoteService.Code)
		test.Equals(t, apiError.Message, ErrorRemoteService.Message)
		test.Equals(t, apiError.Details, "serviceB: serviceA: invalid sql")
	})
}
