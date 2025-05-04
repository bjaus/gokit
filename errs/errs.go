package errs

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/samber/lo"
)

func New(message string) *Error { return &Error{message: message} }

func Wrap(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{cause: err}
}

func Has(err error, codes ...Code) bool {
	if err == nil {
		return false
	}
	e := As(err)
	return slices.Contains(codes, e.Code())
}

func Kind(err error) Code {
	if err == nil {
		return CodeFatal
	}
	e := As(err)
	return lo.Ternary(e == nil, CodeFatal, e.Code())
}

func As(err error) *Error {
	return lo.TernaryF(err == nil,
		func() *Error { return nil },
		func() *Error {
			var e *Error
			return lo.Ternary(errors.As(err, &e), e, Wrap(err))
		},
	)
}

func Is(err error) bool {
	var e *Error
	return lo.TernaryF(err == nil,
		func() bool { return false },
		func() bool { return errors.As(err, &e) },
	)
}

type Code string

const (
	CodeConflict       Code = "conflict"
	CodeFatal          Code = "internal"
	CodeForbidden      Code = "forbidden"
	CodeInvalid        Code = "invalid"
	CodeNotFound       Code = "not-found"
	CodeNotImplemented Code = "not-implemented"
	CodeUnauthorized   Code = "unauthorized"
)

var defaultMessages = map[Code]string{
	CodeConflict:       "conflict",
	CodeFatal:          "internal server error",
	CodeForbidden:      "forbidden",
	CodeInvalid:        "invalid request",
	CodeNotFound:       "resource not found",
	CodeNotImplemented: "not implemented",
	CodeUnauthorized:   "unauthorized",
}

var httpStatusMap = map[Code]int{
	CodeConflict:       http.StatusConflict,
	CodeFatal:          http.StatusInternalServerError,
	CodeForbidden:      http.StatusForbidden,
	CodeInvalid:        http.StatusBadRequest,
	CodeNotFound:       http.StatusNotFound,
	CodeNotImplemented: http.StatusNotImplemented,
	CodeUnauthorized:   http.StatusUnauthorized,
}

type Error struct {
	cause   error
	code    Code
	message string
}

func (e *Error) WithCode(code Code) *Error { e.code = code; return e }

func (e *Error) WithCause(cause error) *Error { e.cause = cause; return e }

func (e *Error) WithMessage(message string) *Error { e.message = message; return e }

func (e *Error) Error() string {
	return lo.Ternary(
		e.cause == nil, e.Message(), fmt.Sprintf("%s: %v", e.Message(), e.cause),
	)
}

func (e *Error) Code() Code {
	return lo.Ternary(e.code == "", CodeFatal, e.code)
}

func (e *Error) Message() string {
	return lo.TernaryF(e.message != "",
		func() string { return e.message },
		func() string {
			if msg, ok := defaultMessages[e.Code()]; ok {
				return msg
			}
			return defaultMessages[CodeFatal]
		},
	)
}

func (e *Error) Unwrap() error { return e.cause }

func (e *Error) status() int {
	if status, ok := httpStatusMap[e.Code()]; ok {
		return status
	}
	return http.StatusInternalServerError
}
